package storagers

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
)

type File struct {
	data  []ShortenedURL
	fname string
}

func (f *File) New(fName string) error {
	f.data = make([]ShortenedURL, 0)
	f.fname = fName
	if err := f.readFromFile(); err != nil {
		return err
	}
	return nil
}

func (f *File) Close() {
}

func (f *File) Set(_ context.Context, index string, value string) error {
	newURL := ShortenedURL{
		UUID:        uuid.New().String(),
		ShortURL:    index,
		OriginalURL: value,
	}
	f.data = append(f.data, newURL)
	err := f.writeValueToFile(newURL)
	if err != nil {
		logger.Log().Error("Error write URL to file", zap.Error(err))
		return err
	}
	logger.Log().Debug("Storage_Set_File", zap.Any("new_url", newURL))
	return nil
}

func (f *File) Get(_ context.Context, url string) (string, error) {
	idx := slices.IndexFunc(f.data, func(c ShortenedURL) bool { return c.ShortURL == url })
	if idx == -1 {
		return "", errors.New("URL not found")
	}
	return f.data[idx].OriginalURL, nil
}

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

func (f *File) writeValueToFile(value ShortenedURL) error {
	file, err := os.OpenFile(f.fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logger.Log().Error("Error create file", zap.Error(err))
		return err
	}
	defer file.Close()

	valByte, err := json.Marshal(&value)
	if err != nil {
		logger.Log().Error("Error marshal json to string", zap.Error(err))
		return err
	}
	valByte = append(valByte, '\n')

	_, err = file.Write(valByte)
	return err

}
