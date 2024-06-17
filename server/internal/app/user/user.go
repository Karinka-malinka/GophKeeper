package user

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID
	Username string
	Password string
}

type UserStore interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, condition map[string]string) (*User, error)
}

type Users struct {
	userStore UserStore
}

func NewUser(userStore UserStore) *Users {
	return &Users{
		userStore: userStore,
	}
}

func (ua *Users) Register(ctx context.Context, user User) (string, error) {

	user.UUID = uuid.New()

	if err := ua.userStore.Create(ctx, user); err != nil {
		return "", err
	}

	accessToken, err := ua.newToken(user, 60, user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (ua *Users) Login(ctx context.Context, user User) (string, error) {

	userInDB, err := ua.userStore.Get(ctx, map[string]string{"login": user.Username})

	if err != nil {
		return "", err
	}

	if !ua.checkHash(user, userInDB.Password) {
		return "", errors.New("401")
	}

	accessToken, err := ua.newToken(*userInDB, 60, user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (ua *Users) checkHash(user User, userHash string) bool {

	check1 := []byte(user.Password)
	check2, err := hex.DecodeString(userHash)

	if err != nil {
		slog.Error("Error in decode user hash. error: " + err.Error())
		return false
	}

	return hmac.Equal(check2, check1)

}
