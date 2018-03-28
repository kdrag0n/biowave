package core

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	l, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
	
	if err != nil {
		panic(err)
	}

	logger = l
}
