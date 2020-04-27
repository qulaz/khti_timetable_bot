package helpers

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

var Logger *zap.SugaredLogger = &zap.SugaredLogger{}

func getLoggerLevel(level string) zap.AtomicLevel {
	var logLevel zapcore.Level
	err := logLevel.UnmarshalText([]byte(level))

	if err != nil {
		log.Fatal(err)
	}

	return zap.NewAtomicLevelAt(logLevel)
}

func dateTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02.01.2006 15:04:05"))
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05"))
}

func InitLogger() {
	cfg := zap.Config{
		Level:            getLoggerLevel(Config.LOG_LEVEL),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     dateTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	if Config.IS_DEBUG {
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeTime = timeEncoder
		cfg.Level = getLoggerLevel("DEBUG")
	} else {
		cfg.Encoding = "json"
		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	logger, err := cfg.Build()

	if err != nil {
		log.Fatal(" Ошибка конфигурации логгера: ", err)
	}

	defer logger.Sync()

	Logger = logger.Sugar()
}
