package logging

import (
	"os"

	log "github.com/sirupsen/logrus"
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
