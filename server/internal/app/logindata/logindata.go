package logindata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// LoginData представляет структуру данных для пары логин/пароль.
type LoginData struct {
	UUID     uuid.UUID
	Created  time.Time
	Login    []byte
	Password []byte
	Meta     []byte
	UserID   string
}

// ILoginDataStore определяет интерфейс хранилища пар логин/пароль
type ILoginDataStore interface {
	Create(ctx context.Context, loginData LoginData) error
	Update(ctx context.Context, loginDataUID string, newpassword []byte) error
	Delete(ctx context.Context, loginDataUID string) error
	GetList(ctx context.Context, userID string) ([]LoginData, error)
}

// LoginDatas представляет структуру для работы с парой логин/пароль
type LoginDatas struct {
	loginDataStore ILoginDataStore
}

// NewLoginData создает новый экземпляр LoginDatas.
func NewLoginData(loginDataStore ILoginDataStore) *LoginDatas {
	return &LoginDatas{
		loginDataStore: loginDataStore,
	}
}

// Add добавляет новую пару логин/пароль.
func (l *LoginDatas) Add(ctx context.Context, loginData LoginData) (*LoginData, error) {

	loginData.UUID = uuid.New()
	loginData.Created = time.Now().UTC()

	if err := l.loginDataStore.Create(ctx, loginData); err != nil {
		return nil, err
	}

	return &loginData, nil
}

// Edit изменяет пару логин/пароль по идентификатору.
func (l *LoginDatas) Edit(ctx context.Context, loginDataUID string, newpassword []byte) error {

	if err := l.loginDataStore.Update(ctx, loginDataUID, newpassword); err != nil {
		return err
	}

	return nil
}

// Delete удаляет пару логин/пароль по идентификатору.
func (l *LoginDatas) Delete(ctx context.Context, loginDataUID string) error {

	if err := l.loginDataStore.Delete(ctx, loginDataUID); err != nil {
		return err
	}

	return nil
}

// GetList возвращает список пар логин/пароль для указанного пользователя.
func (l *LoginDatas) GetList(ctx context.Context, userID string) ([]LoginData, error) {

	list, err := l.loginDataStore.GetList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
