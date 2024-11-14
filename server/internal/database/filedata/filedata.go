package filedata

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/GophKeeper/server/internal/app/filedata"
	"github.com/Masterminds/squirrel"
)

var _ filedata.IFileDataStore = &DataStore{}

// DataStore представляет хранилище данных приватной бинарной информации.
type DataStore struct {
	db *sql.DB
}

// NewFileDataStore создает новый экземпляр хранилища данных приватной бинарной информации.
func NewFileDataStore(db *sql.DB) *DataStore {
	return &DataStore{db: db}
}

// Create создает новую запись приватных бинарных данных в базе данных.
func (d *DataStore) Create(ctx context.Context, data filedata.FileData) error {

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	queryBuilder := qb.Insert("filedata").
		Columns("uuid", "created", "filedata", "namefile", "meta", "user_id").
		Values(data.UUID.String(), data.Created, data.File, data.Name, data.Meta, data.UserID)

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

// Delete удаляет запись приватных бинарной данных из базы данных по идентификатору.
func (d *DataStore) Delete(ctx context.Context, dataUID string) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM filedata WHERE uuid=$1", dataUID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetFile возвращает бинарные данные из базы данных по идентификатору.
func (d *DataStore) GetFile(ctx context.Context, dataUID string) (*filedata.FileData, error) {

	var rows *sql.Rows

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.Select("namefile", "filedata", "meta", "created").
		From("filedata").
		Where(squirrel.Eq{"uui": dataUID}).
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

	var data filedata.FileData

	if rows.Next() {
		if err = rows.Scan(&data.Name, &data.File, &data.Meta, &data.Created); err != nil {
			return nil, err
		}
	}

	return &data, nil
}

// GetList возвращает список приватных бинарной данных для указанного пользователя.
func (d *DataStore) GetList(ctx context.Context, userID string) ([]filedata.FileData, error) {

	var rows *sql.Rows

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.Select("uuid", "namefile", "meta", "created").
		From("filedata").
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

	var datas []filedata.FileData

	for rows.Next() {

		var data filedata.FileData
		if err = rows.Scan(&data.UUID, &data.Name, &data.Meta, &data.Created); err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	return datas, nil
}
