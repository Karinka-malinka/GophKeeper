package filedata

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// FileData представляет структуру данных для бинарной информации
type FileData struct {
	UUID    uuid.UUID
	Created time.Time
	File    []byte
	Meta    []byte
	Name    []byte
	UserID  string
}

// IFileDataStore определяет интерфейс хранилища приватной бинарной информации
type IFileDataStore interface {
	Create(ctx context.Context, data FileData) error
	GetFile(ctx context.Context, dataUID string) (*FileData, error)
	Delete(ctx context.Context, dataUID string) error
	GetList(ctx context.Context, userID string) ([]FileData, error)
}

// FileDatas представляет структуру для работы с приватной бинарной информацией
type FileDatas struct {
	dataStore IFileDataStore
}

// NewFiletData создает новый экземпляр FileDatas.
func NewFiletData(dataStore IFileDataStore) *FileDatas {
	return &FileDatas{
		dataStore: dataStore,
	}
}

// Add добавляет новую приватную бинарную информацию
func (l *FileDatas) Add(ctx context.Context, data FileData) (*FileData, error) {

	data.UUID = uuid.New()
	data.Created = time.Now().UTC()

	if err := l.dataStore.Create(ctx, data); err != nil {
		return nil, err
	}

	return &data, nil
}

// Delete удаляет приватную бинарную информацию
func (l *FileDatas) GetFile(ctx context.Context, dataUID string) (*FileData, error) {

	file, err := l.dataStore.GetFile(ctx, dataUID)

	if err != nil {
		return nil, err
	}

	return file, nil
}

// Delete удаляет приватную бинарную информацию
func (l *FileDatas) Delete(ctx context.Context, dataUID string) error {

	if err := l.dataStore.Delete(ctx, dataUID); err != nil {
		return err
	}

	return nil
}

// GetList возвращает список приватной бинарной информации пользователя.
func (l *FileDatas) GetList(ctx context.Context, userID string) ([]FileData, error) {

	list, err := l.dataStore.GetList(ctx, userID)

	if err != nil {
		return nil, err
	}

	return list, nil
}
