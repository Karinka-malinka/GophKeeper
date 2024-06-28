package textdata

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GophKeeper/server/internal/app/textdata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTextDataStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewTextDataStore(db)

	uid := uuid.New()

	textData := textdata.TextData{
		UUID:    uid,
		Created: time.Now().UTC(),
		Text:    []byte("test text"),
		Meta:    []byte("metadata"),
		UserID:  "456",
	}

	mock.ExpectExec("INSERT INTO textdata").WithArgs(textData.UUID, textData.Created, textData.Text, textData.Meta, textData.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = dataStore.Create(context.Background(), textData)
	assert.NoError(t, err)
}

func TestTextDataStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewTextDataStore(db)

	dataUID := "123"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM textdata").WithArgs(dataUID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dataStore.Delete(context.Background(), dataUID)
	assert.NoError(t, err)
}

func TestTextDataStore_GetList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewTextDataStore(db)

	userID := "456"

	rows := sqlmock.NewRows([]string{"uuid", "textdata", "meta", "created"}).
		AddRow(uuid.New(), "test text", "metadata", time.Now().UTC())

	mock.ExpectQuery("SELECT uuid, textdata, meta, created").WithArgs(userID).WillReturnRows(rows)

	textDatas, err := dataStore.GetList(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, textDatas)
	assert.Len(t, textDatas, 1)
}
