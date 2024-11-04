package uploads

import "fmt"

const (
	DefaultAvatarURL     = "/static/img/KarlMarks.jpg"
	DefaultBackgroundURL = "/static/img/backgroundPicture.png"
)

// JoinFileName восстанавливает имя файла из UUID и расширения (или,
// если UUID пустой, возвращает дефолтное значение)
func JoinFileName(fileUUID string, fileExtension string, defaultValue string) string {
	if fileUUID == "" {
		return defaultValue
	}
	if fileExtension != "" {
		return fmt.Sprintf("%s.%s", fileUUID, fileExtension)
	}
	return fileUUID
}
