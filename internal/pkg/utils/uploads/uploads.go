package uploads

import "fmt"

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
