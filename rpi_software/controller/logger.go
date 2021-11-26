package controller

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(runAsDamon bool) *zap.SugaredLogger {
	var config zap.Config

	if runAsDamon {
		if _, err := os.Stat("/var/log/MKTNE"); os.IsNotExist(err) {
			os.Mkdir("/var/log/MKTNE", 0664)
		}
		config := zap.NewProductionConfig()
		config.OutputPaths = []string{"/var/log/MKTNE/ctrldamon.log"}
		config.Encoding = "console"
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.EncoderConfig.EncodeCaller = nil
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}

	config.EncoderConfig.EncodeName = func(s string, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(fmt.Sprintf("%s", s))
	}
	logger, err := config.Build()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	loggersh := logger.Sugar().Named("MKTNE")
	return loggersh
}
