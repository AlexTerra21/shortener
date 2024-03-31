package storagers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
	"github.com/AlexTerra21/shortener/internal/app/models"
)

// Структура для хранения данных в файле
type File struct {
	data  []ShortenedURL
	fname string
}

// Инициализация хранилища
func (f *File) New(fName string) error {
	f.data = make([]ShortenedURL, 0)
	f.fname = fName
	if err := f.readFromFile(); err != nil {
		return err
	}
	return nil
}

// Закрытие хранилища
func (f *File) Close() {
}

// Добавление данных в хранилище
func (f *File) Set(_ context.Context, index string, value string, userID int) error {
	newURL := ShortenedURL{
		UUID:        userID,
		IdxShortURL: index,
		OriginalURL: value,
		DeletedFlag: false,
	}
	f.data = append(f.data, newURL)
	newURLs := []ShortenedURL{newURL}
	err := f.writeValueToFile(&newURLs)
	if err != nil {
		logger.Log().Error("Error write URL to file", zap.Error(err))
		return err
	}
	logger.Log().Debug("Storage_Set_File", zap.Any("new_url", newURL))
	return nil
}

// Добавление пакета данных в хранилище
func (f *File) BatchSet(_ context.Context, batchValues *[]models.BatchStore, userID int) error {
	newURLs := make([]ShortenedURL, 0)
	for _, url := range *batchValues {
		newURL := ShortenedURL{
			UUID:        userID,
			IdxShortURL: url.IdxShortURL,
			OriginalURL: url.OriginalURL,
			DeletedFlag: false,
		}
		f.data = append(f.data, newURL)
		newURLs = append(newURLs, newURL)
		logger.Log().Debug("Storage_Set_File", zap.Any("new_url", newURL))
	}
	err := f.writeValueToFile(&newURLs)
	if err != nil {
		logger.Log().Error("Error write URL to file", zap.Error(err))
		return err
	}
	return nil
}

// Получение данных из хранилища
func (f *File) Get(_ context.Context, idxURL string) (string, bool, error) {
	idx := slices.IndexFunc(f.data, func(c ShortenedURL) bool { return c.IdxShortURL == idxURL })
	if idx == -1 {
		return "", false, errors.New("URL not found")
	}
	return f.data[idx].OriginalURL, f.data[idx].DeletedFlag, nil
}

// Чтение из файла
func (f *File) readFromFile() error {
	file, err := os.OpenFile(f.fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log().Error("Error open file", zap.Error(err))
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := ShortenedURL{}
		if err := json.Unmarshal(scanner.Bytes(), &val); err != nil {
			logger.Log().Error("Error unmarshal string to json", zap.Error(err))
			return err
		}
		f.data = append(f.data, val)
		logger.Log().Debug("readFromFile", zap.Any("Value", val))
	}
	logger.Log().Sugar().Infof("Shorten URL data restored from %v", f.fname)
	return nil
}

// Запись в файл
func (f *File) writeValueToFile(values *[]ShortenedURL) error {
	file, err := os.OpenFile(f.fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log().Error("Error create file", zap.Error(err))
		return err
	}
	defer file.Close()
	var valByte []byte
	for _, value := range *values {
		val, err := json.Marshal(&value)
		if err != nil {
			logger.Log().Error("Error marshal json to string", zap.Error(err))
			return err
		}
		valByte = append(valByte, val...)
		valByte = append(valByte, '\n')
	}

	_, err = file.Write(valByte)
	return err

}
