package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

// NewLogger создает и настраивает новый экземпляр логгера
func NewLogger(level string) (*Logger, error) {
	logger := logrus.New()

	// Устанавливаем формат логов
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Устанавливаем вывод логов (по умолчанию в консоль)
	logger.SetOutput(os.Stdout)

	// Устанавливаем уровень логирования
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logLevel)

	return &Logger{logger}, nil
}

// ConfigureFileLogger инициализирует логгер для записи в файл
func ConfigureFileLogger(filePath string, level string) (*Logger, error) {
	logger := logrus.New()

	// Открываем или создаем файл для логирования
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	// Устанавливаем файл как выходной поток для логов
	logger.SetOutput(file)

	// Устанавливаем формат логов
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Устанавливаем уровень логирования
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logLevel)

	return &Logger{logger}, nil
}

// Infof записывает информационное сообщение в лог
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Errorf записывает сообщение об ошибке в лог
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Debugf записывает отладочное сообщение в лог
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Warnf записывает предупреждающее сообщение в лог
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}
