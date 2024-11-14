// Package textdata представляет собой пакет для работы с данными текстовых данных.
package textdata

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/GophKeeper/server/internal/app/textdata"
	"github.com/Masterminds/squirrel"
)

var _ textdata.ITextDataStore = &DataStore{}

// DataStore представляет хранилище данных приватной текстовой информации.
type DataStore struct {
	db *sql.DB
}

// NewTextDataStore создает новый экземпляр хранилища данных приватной текстовой информации.
func NewTextDataStore(db *sql.DB) *DataStore {
	return &DataStore{db: db}
}

// Create создает новую запись приватных текстовых данных в базе данных.
func (d *DataStore) Create(ctx context.Context, data textdata.TextData) error {

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	queryBuilder := qb.Insert("textdata").
		Columns("uuid", "created", "textdata", "meta", "user_id").
		Values(data.UUID.String(), data.Created, data.Text, data.Meta, data.UserID)

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

// Delete удаляет запись приватных текстовых данных из базы данных по идентификатору.
func (d *DataStore) Delete(ctx context.Context, dataUID string) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM textdata WHERE uuid=$1", dataUID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetList возвращает список приватных текстовых данных для указанного пользователя.
func (d *DataStore) GetList(ctx context.Context, userID string) ([]textdata.TextData, error) {

	var rows *sql.Rows

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.Select("uuid", "textdata", "meta", "created").
		From("textdata").
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

	var datas []textdata.TextData

	for rows.Next() {

		var data textdata.TextData
		if err = rows.Scan(&data.UUID, &data.Text, &data.Meta, &data.Created); err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	return datas, nil
}
