package textdata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TextData представляет структуру данных для приватной текстовой информации
type TextData struct {
	UUID    uuid.UUID
	Created time.Time
	Text    []byte
	Meta    []byte
	UserID  string
}

// ITextDataStore определяет интерфейс хранилища приватной текстовой информации
type ITextDataStore interface {
	Create(ctx context.Context, textData TextData) error
	Delete(ctx context.Context, textDataUID string) error
	GetList(ctx context.Context, userID string) ([]TextData, error)
}

// TextDatas представляет структуру для работы с приватной текстовой информацией
type TextDatas struct {
	dataStore ITextDataStore
}

// NewTextData создает новый экземпляр TextDatas.
func NewTextData(dataStore ITextDataStore) *TextDatas {
	return &TextDatas{
		dataStore: dataStore,
	}
}

// Add добавляет новую приватную текстовую информацию
func (l *TextDatas) Add(ctx context.Context, data TextData) (*TextData, error) {

	data.UUID = uuid.New()
	data.Created = time.Now().UTC()

	if err := l.dataStore.Create(ctx, data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete удаляет приватную текстовую информацию
func (l *TextDatas) Delete(ctx context.Context, dataUID string) error {

	if err := l.dataStore.Delete(ctx, dataUID); err != nil {
		return err
	}

	return nil
}

// GetList возвращает список приватной текстовой информации пользователя.
func (l *TextDatas) GetList(ctx context.Context, userID string) ([]TextData, error) {

	list, err := l.dataStore.GetList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
