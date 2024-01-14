package storage

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"slices"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	DB    *sql.DB
}

func NewStorage(fname string, dbstr string) (*Storage, error) {

	db, err := sql.Open("pgx", dbstr)
	if err != nil {
		return nil, errors.New("error open database")
	}
	stor := Storage{fname: fname, DB: db}
	if fname != "" { // Отключение чтения из файла
		_ = stor.readFromFile()
	}
	return &stor, nil
}

func (s *Storage) Close() {
	s.DB.Close()
}

func (s *Storage) Set(index string, value string) {
	newURL := shortenedURL{
		UUID:        uuid.New().String(),
		ShortURL:    index,
		OriginalURL: value,
	}
	logger.Log().Debug("Storage_Set", zap.Any("new_url", newURL))
	s.data = append(s.data, newURL)
	if s.fname != "" { // Отключение записи в файл
		err := s.writeValueToFile(newURL)
		if err != nil {
			logger.Log().Error("Error write URL to file", zap.Error(err))
		}
	}
}

func (s *Storage) Get(url string) (string, error) {
	idx := slices.IndexFunc(s.data, func(c shortenedURL) bool { return c.ShortURL == url })
	if idx == -1 {
		return "", errors.New("URL not found")
	}

	return s.data[idx].OriginalURL, nil
}

func (s *Storage) readFromFile() error {
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
