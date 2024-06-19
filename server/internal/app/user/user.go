package user

import (
	"context"
	"errors"

	"github.com/GophKeeper/server/cmd/config"
	"github.com/google/uuid"
)

type User struct {
	UUID     uuid.UUID
	Username string
	Password string
}

type IUserStore interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, login string) (*User, error)
}

type Users struct {
	userStore IUserStore
	cfg       *config.ConfigToken
}

func NewUser(userStore IUserStore, cfg *config.ConfigToken) *Users {
	return &Users{
		userStore: userStore,
		cfg:       cfg,
	}
}

func (ua *Users) Register(ctx context.Context, user User) (string, error) {

	user.UUID = uuid.New()

	if err := ua.userStore.Create(ctx, user); err != nil {
		return "", err
	}

	accessToken, err := ua.newToken(user, ua.cfg.TokenExpiresAt, ua.cfg.SecretKeyForToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (ua *Users) Login(ctx context.Context, user User) (string, error) {

	userInDB, err := ua.userStore.Get(ctx, user.Username)

	if err != nil {
		return "", err
	}

	if user.Password != userInDB.Password {
		return "", errors.New("401")
	}

	accessToken, err := ua.newToken(*userInDB, 60, user.Username)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
