package plterror

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

// InitDevLogger init logger for local development.
// See InitProdLogger for details.
func InitDevLogger() {
	config := newDevConfig()

	logger, err := config.Build()
	if err != nil {
		log.Fatal(err)
	}

	Logger = logger.Sugar()
}

// newDevConfig creates a new dev logger configuration.
func newDevConfig() zap.Config {
	// use default development config as starting base
	config := zap.NewDevelopmentConfig()

	config.Encoding = "console"
	config.EncoderConfig.TimeKey = "time"
	config.EncoderConfig.CallerKey = ""     // disable
	config.EncoderConfig.StacktraceKey = "" // disable
	config.EncoderConfig.EncodeTime = CustomTimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	return config
}

// CustomTimeEncoder serializes a time.Time to a custom formatted string.
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000 Z07"))
}
