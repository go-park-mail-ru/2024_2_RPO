package logging

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

type contextKey string

const (
	RequestIDkey = contextKey("requestID")
)

// SetupLogger настраивает logrus, чтобы он писал в консоль цветом и в файл json-ом
func SetupLogger(jsonFile *os.File) {
	// Настраиваем цвета
	textFormatter := &log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	log.SetFormatter(textFormatter)

	// Устанавливаем вывод в терминал
	log.SetOutput(os.Stdout)

	// Создаем хук для записи JSON-логов в файл
	jsonFormatter := &log.JSONFormatter{}
	log.AddHook(&fileHook{
		Writer:    jsonFile,
		Formatter: jsonFormatter,
	})
}

// fileHook реализует хук для записи логов в файл с определенным форматом
type fileHook struct {
	Writer    *os.File
	Formatter log.Formatter
}

func (hook *fileHook) Levels() []log.Level {
	return log.AllLevels
}

func (hook *fileHook) Fire(entry *log.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

func GetRequestID(ctx context.Context) uint64 {
	requestID, ok := ctx.Value(RequestIDkey).(uint64)
	if !ok {
		return 0
	}
	return requestID
}

func Warn(ctx context.Context, data ...interface{}) {
	logData := append([]interface{}{fmt.Sprintf("rid=%d ", GetRequestID(ctx))}, data...)
	log.Warn(logData...)
}

func Info(ctx context.Context, data ...interface{}) {
	logData := append([]interface{}{fmt.Sprintf("rid=%d ", GetRequestID(ctx))}, data...)
	log.Info(logData...)
}

func Error(ctx context.Context, data ...interface{}) {
	logData := append([]interface{}{fmt.Sprintf("rid=%d ", GetRequestID(ctx))}, data...)
	log.Error(logData...)
}

func Debug(ctx context.Context, data ...interface{}) {
	logData := append([]interface{}{fmt.Sprintf("rid=%d ", GetRequestID(ctx))}, data...)
	log.Info(logData...)
}

func Warnf(ctx context.Context, format string, data ...interface{}) {
	ridInfo := fmt.Sprintf("rid=%d ", GetRequestID(ctx))
	log.Warnf(ridInfo+format, data...)
}

func Infof(ctx context.Context, format string, data ...interface{}) {
	ridInfo := fmt.Sprintf("rid=%d ", GetRequestID(ctx))
	log.Infof(ridInfo+format, data...)
}

func Errorf(ctx context.Context, format string, data ...interface{}) {
	ridInfo := fmt.Sprintf("rid=%d ", GetRequestID(ctx))
	log.Errorf(ridInfo+format, data...)
}

func Debugf(ctx context.Context, format string, data ...interface{}) {
	ridInfo := fmt.Sprintf("rid=%d ", GetRequestID(ctx))
	log.Debugf(ridInfo+format, data...)
}
