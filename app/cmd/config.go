package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Site struct {
	BaseUrl string
}

type Config struct {
	Debug          bool
	Port           int
	DbFilePath     string
	DbSaveInterval string
}

type GlobalConfig struct {
	Db_path string
	Host    string
	Go_env  string
	Port    int
}

var site = Site{}
var zapLogger *zap.Logger

func (config *Config) init() {
	var goenv GlobalConfig
	envconfig.Process("server", &goenv)
	if goenv.Go_env == "development" {
		config.Debug = true
	} else {
		config.Debug = false
	}
	zapLogger, _ = configureLogger(config.Debug)

	config.Port = goenv.Port
	if config.Debug {
		zapLogger.Info("Debug mode enabled")
	}

	config.Port = goenv.Port
	zapLogger.Debug("Port set", zap.Int("port", config.Port))
	site.BaseUrl = "http://" + goenv.Host + ":" + fmt.Sprint(goenv.Port)
	zapLogger.Debug("Site base URL set", zap.String("base_url", site.BaseUrl))
}

func configureLogger(debug bool) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	if debug {
		level = zapcore.DebugLevel
	}

	zapConfig := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encoderConfig,
	}
	return zapConfig.Build()
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05"))
}
