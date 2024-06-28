package bankcard

import (
	"context"
	"time"
)

// BankCardData представляет структуру данных о банковской карте
type BankCardData struct {
	Number  []byte
	Created time.Time
	Term    []byte
	CCV     []byte
	Meta    []byte
	UserID  string
}

// IBankCardDataStore определяет интерфейс хранилища данных о банковской карте
type IBankCardDataStore interface {
	Create(ctx context.Context, data BankCardData) error
	Delete(ctx context.Context, dataUID string) error
	GetList(ctx context.Context, userID string) ([]BankCardData, error)
}

// BankCardDatas представляет структуру для работы с данными о банковской карте
type BankCardDatas struct {
	dataStore IBankCardDataStore
}

// NewBankCardData создает новый экземпляр BankCardDatas.
func NewBankCardData(dataStore IBankCardDataStore) *BankCardDatas {
	return &BankCardDatas{
		dataStore: dataStore,
	}
}

// Add добавляет новые данные о банковской карте
func (l *BankCardDatas) Add(ctx context.Context, data BankCardData) (*BankCardData, error) {

	data.Created = time.Now().UTC()

	if err := l.dataStore.Create(ctx, data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete удаляет данные о банковской карте
func (l *BankCardDatas) Delete(ctx context.Context, dataUID string) error {

	if err := l.dataStore.Delete(ctx, dataUID); err != nil {
		return err
	}

	return nil
}

// GetList возвращает список данных о банковских картах пользователя.
func (l *BankCardDatas) GetList(ctx context.Context, userID string) ([]BankCardData, error) {

	list, err := l.dataStore.GetList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
