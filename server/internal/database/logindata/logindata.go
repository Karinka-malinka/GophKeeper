package logindata

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/GophKeeper/server/internal/app/logindata"
	"github.com/Masterminds/squirrel"
)

var _ logindata.ILoginDataStore = &LoginDataStore{}

type LoginDataStore struct {
	db *sql.DB
}

func NewLoginDataStore(db *sql.DB) *LoginDataStore {
	return &LoginDataStore{db: db}
}

func (d *LoginDataStore) Create(ctx context.Context, loginData logindata.LoginData) error {

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	queryBuilder := qb.Insert("logindata").
		Columns("uuid", "created", "login", "password", "meta", "user_id").
		Values(loginData.UUID.String(), loginData.Created, loginData.Login, loginData.Password, loginData.Meta, loginData.UserID)

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

func (d *LoginDataStore) Update(ctx context.Context, loginDataUID string, newpassword []byte) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE logindata SET password=$1 WHERE uuid=$2", newpassword, loginDataUID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *LoginDataStore) Delete(ctx context.Context, loginDataUID string) error {

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM logindata WHERE uuid=$1", loginDataUID)

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *LoginDataStore) GetList(ctx context.Context, userID string) ([]logindata.LoginData, error) {

	var rows *sql.Rows

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.Select("uuid", "login", "password", "meta", "created").
		From("logindata").
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

	var logindatas []logindata.LoginData

	for rows.Next() {

		var loginData logindata.LoginData
		if err = rows.Scan(&loginData.UUID, &loginData.Login, &loginData.Password, &loginData.Meta, &loginData.Created); err != nil {
			return nil, err
		}

		logindatas = append(logindatas, loginData)
	}

	return logindatas, nil
}
