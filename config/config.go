package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// LogConfig структура для конфигурации логирования
type LogConfig struct {
	Level      string `yaml:"level"`
	File       string `yaml:"file"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	Compress   bool   `yaml:"compress"`
}

// DatabaseConfig структура для конфигурации базы данных
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// ServerConfig структура для конфигурации сервера
type ServerConfig struct {
	Port string `yaml:"port"`
}

// AppConfig структура для общей конфигурации приложения
type AppConfig struct {
	Logging  LogConfig      `yaml:"logging"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// LoadConfig загружает конфигурацию приложения из YAML файла
func LoadConfig(configPath string) (*AppConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Если путь к файлу логов указан как относительный, преобразуем его в абсолютный
	if config.Logging.File != "" && !filepath.IsAbs(config.Logging.File) {
		config.Logging.File = filepath.Join(filepath.Dir(configPath), config.Logging.File)
	}

	return &config, nil
}

// ConfigureLogger настраивает логирование на основе конфигурационного файла
func ConfigureLogger(configPath string) {
	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации логирования: %v", err)
	}

	// Установка уровня логирования
	level, err := logrus.ParseLevel(config.Logging.Level)
	if err != nil {
		log.Fatalf("Ошибка при установке уровня логирования: %v", err)
	}
	logrus.SetLevel(level)

	// Обработка файла логирования
	if config.Logging.File != "" {
		// Если путь относительный, делаем его абсолютным на основе configPath
		if !filepath.IsAbs(config.Logging.File) {
			config.Logging.File = filepath.Join(filepath.Dir(configPath), config.Logging.File)
		}

		// Создание директории, если ее нет
		dir := filepath.Dir(config.Logging.File)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Ошибка при создании директории для лог-файла: %v", err)
		}

		// Открытие файла для логирования
		file, err := os.OpenFile(config.Logging.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Ошибка при открытии файла для логирования: %v", err)
		}
		logrus.SetOutput(file)
	} else {
		logrus.Warn("Не указан файл для логирования в конфигурации, логирование будет осуществляться в консоль")
	}
}

// GetDatabaseURL возвращает строку подключения к базе данных
func GetDatabaseURL(config *AppConfig) string {
	return "host=" + config.Database.Host +
		" port=" + config.Database.Port +
		" user=" + config.Database.User +
		" password=" + config.Database.Password +
		" dbname=" + config.Database.DBName +
		" sslmode=" + config.Database.SSLMode
}

// LoadAppConfig загружает конфигурацию приложения и инициализирует настройки
func LoadAppConfig(configPath string) (*AppConfig, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	// Конфигурация логирования
	ConfigureLogger(configPath)

	return config, nil
}
