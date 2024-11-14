package logindata

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestLoginDataStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	loginDataStore := NewLoginDataStore(db)

	uid := uuid.New()

	loginData := logindata.LoginData{
		UUID:     uid,
		Created:  time.Now().UTC(),
		Login:    []byte("testuser"),
		Password: []byte("testpass"),
		Meta:     []byte("metadata"),
		UserID:   "456",
	}

	mock.ExpectExec("INSERT INTO logindata").WithArgs(loginData.UUID, loginData.Created, loginData.Login, loginData.Password, loginData.Meta, loginData.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = loginDataStore.Create(context.Background(), loginData)
	assert.NoError(t, err)
}

func TestLoginDataStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	loginDataStore := NewLoginDataStore(db)

	loginDataUID := "123"
	newPassword := []byte("newpass")

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE logindata").WithArgs(newPassword, loginDataUID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = loginDataStore.Update(context.Background(), loginDataUID, newPassword)
	assert.NoError(t, err)
}

func TestLoginDataStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	loginDataStore := NewLoginDataStore(db)

	loginDataUID := "123"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM logindata").WithArgs(loginDataUID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = loginDataStore.Delete(context.Background(), loginDataUID)
	assert.NoError(t, err)
}

func TestLoginDataStore_GetList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	loginDataStore := NewLoginDataStore(db)

	userID := "456"

	rows := sqlmock.NewRows([]string{"uuid", "login", "password", "meta", "created"}).
		AddRow(uuid.New(), "testuser", "testpass", "metadata", time.Now().UTC())

	mock.ExpectQuery("SELECT uuid, login, password, meta, created").WithArgs(userID).WillReturnRows(rows)

	loginDatas, err := loginDataStore.GetList(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, loginDatas)
	assert.Len(t, loginDatas, 1)
}
