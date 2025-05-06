package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func Init() {
	log = logrus.New()

	// Уровень логирования (можно сделать настраиваемым через переменные окружения)
	log.SetLevel(logrus.DebugLevel)

	// Формат логов
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableQuote:     true,
		DisableTimestamp: false,
	})

	// Вывод (stdout)
	log.SetOutput(os.Stdout)
}

// Утилиты для использования логгера из других пакетов

func Info(args ...interface{}) {
	log.Info(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}
