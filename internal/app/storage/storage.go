package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/AlexTerra21/shortener/internal/app/logger"
)

type shortenedURL struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Storage struct {
	fname string
	data  []shortenedURL
}

func NewStorage(fname string) *Storage {
	stor := Storage{fname: fname}
	_ = stor.readFromFile()
	return &stor
}

func (s *Storage) Close() {
}

func (s *Storage) Set(index string, value string) {
	new_url := shortenedURL{
		UUID:        uuid.New().String(),
		ShortURL:    index,
		OriginalURL: value,
	}
	logger.Log().Debug("Storage_Set", zap.Any("new_url", new_url))
	s.data = append(s.data, new_url)
	err := s.writeValueToFile(new_url)
	if err != nil {
		logger.Log().Error("Error write URL to file", zap.Error(err))
	}
}

func (s *Storage) Get(url string) (string, error) {
	idx := slices.IndexFunc(s.data, func(c shortenedURL) bool { return c.ShortURL == url })
	if idx == -1 {
		return "", errors.New("")
	}

	return s.data[idx].OriginalURL, nil
}

func (s *Storage) readFromFile() error {
	if s.fname == "" {
		return nil
	}

	file, err := os.OpenFile(s.fname, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Log().Error("Error open file", zap.Error(err))
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		val := shortenedURL{}
		if err := json.Unmarshal(scanner.Bytes(), &val); err != nil {
			logger.Log().Error("Error unmarshal string to json", zap.Error(err))
			return err
		}
		s.data = append(s.data, val)
		logger.Log().Debug("readFromFile", zap.Any("Value", val))
	}
	logger.Log().Sugar().Infof("Shorten URL data restored from %v", s.fname)
	return nil
}

func (s *Storage) writeValueToFile(value shortenedURL) error {
	if s.fname == "" {
		return nil
	}
	file, err := os.OpenFile(s.fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
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
