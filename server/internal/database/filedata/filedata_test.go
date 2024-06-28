package filedata

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GophKeeper/server/internal/app/filedata"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFileDataStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewFileDataStore(db)

	fileData := filedata.FileData{
		UUID:    uuid.New(),
		Created: time.Now().UTC(),
		File:    []byte{65, 66, 67, 68, 69},
		Name:    []byte("testfile.txt"),
		Meta:    []byte("metadata"),
		UserID:  "456",
	}

	mock.ExpectExec("INSERT INTO filedata").WithArgs(fileData.UUID, fileData.Created, fileData.File, fileData.Name, fileData.Meta, fileData.UserID).WillReturnResult(sqlmock.NewResult(1, 1))

	err = dataStore.Create(context.Background(), fileData)
	assert.NoError(t, err)
}

func TestFileDataStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewFileDataStore(db)

	dataUID := "123"

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM filedata").WithArgs(dataUID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = dataStore.Delete(context.Background(), dataUID)
	assert.NoError(t, err)
}

func TestFileDataStore_GetFile(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewFileDataStore(db)

	dataUID := "123"

	rows := sqlmock.NewRows([]string{"namefile", "filedata", "meta", "created"}).
		AddRow([]byte("testfile.txt"), []byte{65, 66, 67, 68, 69}, "metadata", time.Now().UTC())

	mock.ExpectQuery("SELECT namefile, filedata, meta, created").WithArgs(dataUID).WillReturnRows(rows)

	fileData, err := dataStore.GetFile(context.Background(), dataUID)
	assert.NoError(t, err)
	assert.NotNil(t, fileData)
}

func TestFileDataStore_GetList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataStore := NewFileDataStore(db)

	userID := "456"

	rows := sqlmock.NewRows([]string{"uuid", "namefile", "meta", "created"}).
		AddRow(uuid.New(), "testfile.txt", "metadata", time.Now().UTC())

	mock.ExpectQuery("SELECT uuid, namefile, meta, created").WithArgs(userID).WillReturnRows(rows)

	fileDatas, err := dataStore.GetList(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, fileDatas)
	assert.Len(t, fileDatas, 1)
}
