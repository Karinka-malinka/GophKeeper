package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

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
	user.Password = hex.EncodeToString(ua.writeHash(user.Username, user.Password))

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

	check1 := ua.writeHash(user.Username, user.Password)
	check2, err := hex.DecodeString(userHash)

	if err != nil {
		log.Printf("Error in decode user hash. error: %v\n", err)
	}

	return hmac.Equal(check2, check1)

}

func (ua *Users) writeHash(username string, password string) []byte {

	hash := hmac.New(sha256.New, []byte("ua.Cfg.SecretKeyForHashingPassword"))
	hash.Write([]byte(fmt.Sprintf("%s:%s:%s", username, password, "ua.Cfg.SecretKeyForHashingPassword")))

	return hash.Sum(nil)
}
