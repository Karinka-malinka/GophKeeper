package bankcard

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/GophKeeper/server/internal/app/bankcard"
	"github.com/Masterminds/squirrel"
)

var _ bankcard.IBankCardDataStore = &DataStore{}

// DataStore представляет хранилище банковских данных.
type DataStore struct {
	db *sql.DB
}

// NewBankCardDataStore создает новый экземпляр хранилища банковских данных.
func NewBankCardDataStore(db *sql.DB) *DataStore {
	return &DataStore{db: db}
}

// Create создает новую запись банковских данных в базе данных.
func (d *DataStore) Create(ctx context.Context, data bankcard.BankCardData) error {

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	queryBuilder := qb.Insert("bankcard").
		Columns("numbercard", "created", "term", "ccv", "meta", "user_id").
		Values(data.Number, data.Created, data.Term, data.CCV, data.Meta, data.UserID)

	// Получаем SQL запрос и параметры
	sqlString, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = d.db.ExecContext(ctx, sqlString, args...)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

// Delete удаляет запись банковских данных из базы данных по идентификатору.
func (d *DataStore) Delete(ctx context.Context, dataUID []byte) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM bankcard WHERE numbercard=$1", dataUID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetList возвращает список банковских данных для указанного пользователя.
func (d *DataStore) GetList(ctx context.Context, userID string) ([]bankcard.BankCardData, error) {

	var rows *sql.Rows

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.Select("numbercard", "term", "ccv", "meta", "created").
		From("bankcard").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err = d.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	defer rows.Close()

	var datas []bankcard.BankCardData

	for rows.Next() {

		var data bankcard.BankCardData
		if err = rows.Scan(&data.Number, &data.Term, &data.CCV, &data.Meta, &data.Created); err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	return datas, nil
}
