package user

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/GophKeeper/server/internal/app/user"
	"github.com/GophKeeper/server/internal/database"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var _ user.IUserStore = &UserStore{}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {

	return &UserStore{db: db}
}

func (d *UserStore) Create(ctx context.Context, user user.User) error {

	_, err := d.db.ExecContext(ctx, "INSERT INTO users (uuid, login, hash_pass) VALUES($1,$2,$3)", user.UUID.String(), user.Username, user.Password)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			return database.NewErrorConflict(err)
		}

		slog.Error(err.Error())
		return err
	}

	return nil
}

func (d *UserStore) Get(ctx context.Context, login string) (*user.User, error) {
	var rows *sql.Rows

	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	query, args, err := qb.Select("uuid, login, hash_pass").
		From("users").
		Where(squirrel.Eq{"login": login}).
		ToSql()

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	rows, err = d.db.QueryContext(ctx, query, args...)

	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if rows.Err() != nil {
		slog.Error(rows.Err().Error())
		return nil, rows.Err()
	}

	defer rows.Close()

	var user user.User

	for rows.Next() {

		if err = rows.Scan(&user.UUID, &user.Username, &user.Password); err != nil {
			return nil, errors.New("401")
		}
	}

	return &user, nil
}
