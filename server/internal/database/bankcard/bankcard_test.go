package bankcard

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GophKeeper/server/internal/app/bankcard"
	"github.com/stretchr/testify/assert"
)

func TestCreateBankCardData(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dataStore := NewBankCardDataStore(db)

	data := bankcard.BankCardData{
		Number:  []byte("1234567890123456"),
		Created: time.Now().UTC(),
		Term:    []byte("12/25"),
		CCV:     []byte("123"),
		Meta:    []byte("Test meta"),
		UserID:  "456",
	}

	mock.ExpectExec("INSERT INTO bankcard").WillReturnResult(sqlmock.NewResult(1, 1))

	err := dataStore.Create(context.Background(), data)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBankCardData(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dataStore := NewBankCardDataStore(db)

	dataUID := []byte("1234567890123456")

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM bankcard").WithArgs(dataUID).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := dataStore.Delete(context.Background(), dataUID)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetListBankCardData(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	dataStore := NewBankCardDataStore(db)

	userID := "user123"

	rows := sqlmock.NewRows([]string{"numbercard", "term", "ccv", "meta", "created"}).
		AddRow("1234567890123456", "12/25", "123", "Test meta", time.Now().UTC())

	mock.ExpectQuery("SELECT numbercard, term, ccv, meta, created FROM bankcard").WithArgs(userID).WillReturnRows(rows)

	dataList, err := dataStore.GetList(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, dataList)
	assert.NoError(t, mock.ExpectationsWereMet())
}
