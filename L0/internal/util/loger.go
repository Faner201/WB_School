package util

import (
	"github.com/rs/zerolog"
)

type Logger struct {
	Logger *zerolog.Logger
}

func (l *Logger) InitLogger(loggerLevel string) error {
	level, err := zerolog.ParseLevel(loggerLevel)
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(level)
	return nil
}
