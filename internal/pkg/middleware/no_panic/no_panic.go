package no_panic

import (
	"RPO_back/internal/pkg/utils/responses"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
)

func PanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Паника для прерывания запроса обрабатывается в mux'e
				if err == http.ErrAbortHandler {
					log.Warn("Abort connection")
					panic(err)
				}
				log.Error("Panic: ", err)
				log.Error("Debug stack: ", prettifyStack(debug.Stack(), err))
				responses.DoBadResponse(w, http.StatusInternalServerError, "internal error")
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Остальная часть файла взята из библиотеки chi с изменениями
// https://github.com/go-chi/chi
// MIT License

type prettyStack struct {
}

func prettifyStack(debugStack []byte, err interface{}) string {
	s := prettyStack{}
	out, err := s.parse(debugStack, err)
	if err == nil {
		return out
	}
	return "Can't print stack to log"
}

func (s prettyStack) parse(debugStack []byte, rvr interface{}) (string, error) {
	var err error
	buf := &strings.Builder{}

	buf.WriteString("\n")
	buf.WriteString(" panic: ")
	buf.WriteString(fmt.Sprintf("%v", rvr))
	buf.WriteString("\n \n")

	// Обработка информации стека отладки
	stack := strings.Split(string(debugStack), "\n")
	lines := []string{}

	// Поиск строки паники, так как может быть вложенная паника
	for i := len(stack) - 1; i > 0; i-- {
		lines = append(lines, stack[i])
		if strings.HasPrefix(stack[i], "panic(") {
			lines = lines[0 : len(lines)-2] // удаление шаблонного текста
			break
		}
	}

	// Реверсирование списка строк
	for i := len(lines)/2 - 1; i >= 0; i-- {
		opp := len(lines) - 1 - i
		lines[i], lines[opp] = lines[opp], lines[i]
	}

	// Добавление украшений к каждой строке
	for i, line := range lines {
		lines[i], err = s.decorateLine(line, i)
		if err != nil {
			return "", err
		}
	}

	for _, l := range lines {
		buf.WriteString(l)
	}
	return buf.String(), nil
}

func (s prettyStack) decorateLine(line string, num int) (string, error) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "\t") || strings.Contains(line, ".go:") {
		return s.decorateSourceLine(line, num)
	}
	if strings.HasSuffix(line, ")") {
		return s.decorateFuncCallLine(line, num)
	}
	if strings.HasPrefix(line, "\t") {
		return strings.Replace(line, "\t", "      ", 1), nil
	}
	return fmt.Sprintf("    %s\n", line), nil
}

func (s prettyStack) decorateFuncCallLine(line string, num int) (string, error) {
	idx := strings.LastIndex(line, "(")
	if idx < 0 {
		return "", errors.New("не является строкой вызова функции")
	}

	var result strings.Builder
	pkg := line[0:idx]
	method := ""

	if idx := strings.LastIndex(pkg, string(os.PathSeparator)); idx < 0 {
		if idx := strings.Index(pkg, "."); idx > 0 {
			method = pkg[idx:]
			pkg = pkg[0:idx]
		}
	} else {
		method = pkg[idx+1:]
		pkg = pkg[0 : idx+1]
		if idx := strings.Index(method, "."); idx > 0 {
			pkg += method[0:idx]
			method = method[idx:]
		}
	}

	if num == 0 {
		result.WriteString(" -> ")
	} else {
		result.WriteString("    ")
	}
	result.WriteString(pkg)
	result.WriteString(method)
	result.WriteString("\n")
	return result.String(), nil
}

func (s prettyStack) decorateSourceLine(line string, num int) (string, error) {
	idx := strings.LastIndex(line, ".go:")
	if idx < 0 {
		return "", errors.New("не является исходной строкой")
	}

	var result strings.Builder
	path := line[0 : idx+3]
	lineno := line[idx+3:]

	idx = strings.LastIndex(path, string(os.PathSeparator))
	dir := path[0 : idx+1]
	file := path[idx+1:]

	idx = strings.Index(lineno, " ")
	if idx > 0 {
		lineno = lineno[0:idx]
	}

	if num == 1 {
		result.WriteString(" ->   ")
	} else {
		result.WriteString("      ")
	}
	result.WriteString(dir)
	result.WriteString(file)
	result.WriteString(lineno)
	result.WriteString("\n")

	return result.String(), nil
}
