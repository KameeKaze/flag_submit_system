package utils

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func init() {
	Logging()
}

var Logger *zap.Logger

//logging to log.json and terminal
//create logger
func Logging() {
	// the log file
	logFile, _ := os.OpenFile("fss.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// create config
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder // set iso time format

	// create encoder for both terminal and logfile
	fileEncoder := zapcore.NewJSONEncoder(config)
	consoleEncoder := zapcore.NewConsoleEncoder(config)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), zap.DebugLevel),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
	)
	Logger = zap.New(core, zap.AddCaller())

}

func GenerateToken() string {
	return uuid.New().String()
}

func MsgParser(input string) string {
	return `{"msg": "` + input + `"}`
}