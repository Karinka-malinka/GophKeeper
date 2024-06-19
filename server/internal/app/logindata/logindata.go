package logindata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LoginData struct {
	UUID     uuid.UUID
	Created  time.Time
	Login    []byte
	Password []byte
	Meta     []byte
	UserID   string
}

type ILoginDataStore interface {
	Create(ctx context.Context, loginData LoginData) error
	Update(ctx context.Context, loginDataUID, newpassword string) error
	Delete(ctx context.Context, loginDataUID string) error
	GetList(ctx context.Context, userID string) ([]LoginData, error)
}

type LoginDatas struct {
	loginDataStore ILoginDataStore
}

func NewUser(loginDataStore ILoginDataStore) *LoginDatas {
	return &LoginDatas{
		loginDataStore: loginDataStore,
	}
}

func (l *LoginDatas) Add(ctx context.Context, loginData LoginData) error {

	loginData.UUID = uuid.New()
	loginData.Created = time.Now().UTC()

	if err := l.loginDataStore.Create(ctx, loginData); err != nil {
		return err
	}

	return nil
}

func (l *LoginDatas) Edit(ctx context.Context, loginDataUID, newpassword string) error {

	if err := l.loginDataStore.Update(ctx, loginDataUID, newpassword); err != nil {
		return err
	}

	return nil
}

func (l *LoginDatas) Delete(ctx context.Context, loginDataUID string) error {

	if err := l.loginDataStore.Delete(ctx, loginDataUID); err != nil {
		return err
	}

	return nil
}

func (l *LoginDatas) GetList(ctx context.Context, userID string) ([]LoginData, error) {

	list, err := l.loginDataStore.GetList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
