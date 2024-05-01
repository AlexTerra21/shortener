package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/AlexTerra21/shortener/internal/app/async"
	"github.com/AlexTerra21/shortener/internal/app/storage"
)

// Структура для хранения конфигурации приложения
type Config struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	LogLevel        string `json:"log_level"`
	FileStoragePath string `json:"file_storage_path"`
	DBConnectString string `json:"db_connect_string"`
	EnableHTTPS     bool   `json:"enable_https"`
	ConfigPath      string
	Storage         *storage.Storage
	DelQueue        *async.Async
}

// Инициализация конфигурации
func NewConfig() (*Config, error) {
	config := &Config{}
	// Обработка флагов командной строки и занесение в параметры конфигурации
	flagServerAddress := flag.String("a", "", "address and port to run server")
	flagBaseURL := flag.String("b", "", "address and port to return")
	flagLogLevel := flag.String("l", "", "log level")
	flagFileStoragePath := flag.String("f", "", "file name for url save")
	flagDBConnectString := flag.String("d", "", "db connection string")
	flagEnableHTTPS := flag.String("s", "", "enable HTTPS")
	flagConfigPath := flag.String("c", "", "config path")
	flag.Parse()

	config.ConfigPath = *flagConfigPath

	configFromFile, err := config.ReadConFile(config.ConfigPath)
	if err != nil {
		return &Config{}, err
	}

	config.BaseURL = priorityString(os.Getenv("BASE_URL"), *flagBaseURL, configFromFile.BaseURL, "http://localhost:8080")
	config.DBConnectString = priorityString(os.Getenv("DATABASE_DSN"), *flagDBConnectString, configFromFile.DBConnectString)
	config.FileStoragePath = priorityString(os.Getenv("FILE_STORAGE_PATH"), *flagFileStoragePath, configFromFile.FileStoragePath)
	config.LogLevel = priorityString(os.Getenv("LOG_LEVEL"), *flagLogLevel, configFromFile.LogLevel, "info")
	config.ServerAddress = priorityString(os.Getenv("SERVER_ADDRESS"), *flagServerAddress, configFromFile.ServerAddress, ":8080")
	enableHTTPS := priorityString(os.Getenv("ENABLE_HTTPS"), *flagEnableHTTPS, strconv.FormatBool(configFromFile.EnableHTTPS), "false")

	if boolValue, err := strconv.ParseBool(enableHTTPS); err == nil {
		config.EnableHTTPS = boolValue
	} else {
		config.EnableHTTPS = false
	}

	return config, nil
}

// Инициализация хранилища
func (c *Config) InitStorage() (err error) {
	// logger.Log().Info(c.dbConnectString)
	c.Storage, err = storage.NewStorage(c.FileStoragePath, c.DBConnectString)
	return
}

// Инициализация процесса для удаления записей из БД
func (c *Config) InitAsync() {
	c.DelQueue = async.NewAsync(c.Storage)
}

// Присваивание параметра конфигурации serverAddress
func (c *Config) SetServerAddress(s string) {
	c.ServerAddress = s
}

// Присваивание параметра конфигурации baseURL
func (c *Config) SetBaseURL(s string) {
	c.BaseURL = s
}

// Печать основных параметров конфигурации
func (c *Config) Print() {
	fmt.Printf("Server address: %s\n", c.ServerAddress)
	fmt.Printf("Base URL: %s\n", c.BaseURL)
	fmt.Printf("Log level: %s\n", c.LogLevel)
	fmt.Printf("DB connect string: %s\n", c.DBConnectString)
	fmt.Printf("File storage path: %s\n", c.FileStoragePath)
	fmt.Printf("Enable HTTPS: %v\n", c.EnableHTTPS)

}

// Выбор первой не пустой строки по порядку приоритета
func priorityString(vars ...string) string {
	for _, v := range vars {
		if v != "" {
			return v
		}
	}
	return ""
}

// Чтение файла конфигурации
func (c *Config) ReadConFile(path string) (Config, error) {
	if path == "" {
		return Config{}, nil
	}

	fileConfig, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("can not read config from - %s", path)
		}
		return Config{}, err
	}

	config := Config{}
	err = json.Unmarshal(fileConfig, &config)
	return config, err
}
