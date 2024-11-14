package user

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/GophKeeper/server/internal/app/user"
	"github.com/GophKeeper/server/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestUserStore_Create(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStore := NewUserStore(db)

	user := user.User{
		UUID:     uuid.New(),
		Username: "testuser",
		Password: "testpass",
	}

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO users").WithArgs(user.UUID, user.Username, user.Password).WillReturnResult(sqlmock.NewResult(1, 1))

		err = userStore.Create(context.Background(), user)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {

		mock.ExpectExec("INSERT INTO users").WithArgs(user.UUID, user.Username, user.Password).WillReturnError(errors.New("err"))

		err = userStore.Create(context.Background(), user)
		assert.Error(t, err)
	})

	t.Run("Conflict", func(t *testing.T) {

		pgErr := &pgconn.PgError{
			Code: "23505", // Код ошибки нарушения уникального ограничения
		}

		mock.ExpectExec("INSERT INTO users").WithArgs(user.UUID, user.Username, user.Password).WillReturnError(pgErr)

		err = userStore.Create(context.Background(), user)
		assert.Error(t, err)
		assert.Equal(t, err, &database.ConflictError{Err: pgErr})
	})
}

func TestUserStore_GetUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userStore := NewUserStore(db)
	uid := uuid.New()
	rows := sqlmock.NewRows([]string{"uuid", "login", "hash_pass"}).
		AddRow(uid.String(), "testuser", "testpass")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery("SELECT uuid, login, hash_pass").WithArgs("testuser").WillReturnRows(rows)

		user, err := userStore.GetUser(context.Background(), "testuser")
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, uid, user.UUID)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "testpass", user.Password)
	})

	t.Run("Error", func(t *testing.T) {
		mock.ExpectQuery("SELECT uuid, login, hash_pass").WithArgs("testuser").WillReturnError(errors.New("Error"))

		_, err := userStore.GetUser(context.Background(), "testuser")
		assert.Error(t, err)
	})
}
