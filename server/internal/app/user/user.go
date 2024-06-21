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
	Token    string
}

type IUserStore interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, login string) (*User, error)
}

type Users struct {
	userStore IUserStore
	Cfg       *config.ConfigToken
}

func NewUser(userStore IUserStore, cfg *config.ConfigToken) *Users {
	return &Users{
		userStore: userStore,
		Cfg:       cfg,
	}
}

func (ua *Users) Register(ctx context.Context, user User) (*User, error) {

	user.UUID = uuid.New()

	if err := ua.userStore.Create(ctx, user); err != nil {
		return nil, err
	}

	accessToken, err := ua.newToken(user, ua.Cfg.TokenExpiresAt, ua.Cfg.SecretKeyForToken)
	if err != nil {
		return nil, err
	}

	user.Token = accessToken

	return &user, nil
}

func (ua *Users) Login(ctx context.Context, user User) (*User, error) {

	userInDB, err := ua.userStore.Get(ctx, user.Username)

	if err != nil {
		return nil, err
	}

	if user.Password != userInDB.Password {
		return nil, errors.New("401")
	}

	accessToken, err := ua.newToken(*userInDB, ua.Cfg.TokenExpiresAt, ua.Cfg.SecretKeyForToken)
	if err != nil {
		return nil, err
	}

	userInDB.Token = accessToken

	return userInDB, nil
}
