package uploads

import (
	"fmt"
	"os"
	"strings"
)

const (
	DefaultAvatarURL     = "/static/img/KarlMarks.jpg"
	DefaultBackgroundURL = "/static/img/backgroundPicture.png"
)

// JoinFileName восстанавливает URL файла из UUID и расширения (или,
// если UUID пустой, возвращает дефолтное значение)
func JoinFileURL(fileUUID string, fileExtension string, defaultValue string) string {
	urlPrefix := os.Getenv("USER_UPLOADS_URL")
	if fileUUID == "" {
		return urlPrefix + defaultValue
	}
	if fileExtension != "" {
		return urlPrefix + fmt.Sprintf("%s.%s", fileUUID, fileExtension)
	}
	return urlPrefix + fileUUID
}

// ExtractFileExtension извлекает из имени файла его расширение (возвращает "", если файл без расширения)
func ExtractFileExtension(fileName string) string {
	lastDotIndex := strings.LastIndex(fileName, ".")
	if lastDotIndex == -1 || lastDotIndex == len(fileName)-1 {
		return ""
	}
	return fileName[lastDotIndex+1:]
}

// JoinFilePath восстанавливает имя файла из UUID
func JoinFilePath(fileUUID string, fileExtension string) string {
	if fileExtension != "" {
		return fmt.Sprintf("%s.%s", fileUUID, fileExtension)
	}
	return fileUUID
}
