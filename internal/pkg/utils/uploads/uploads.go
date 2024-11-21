package uploads

import (
	"RPO_back/internal/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
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
		return defaultValue
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

func CompareFiles(fileNames []string, newFile []byte) (fileUUID string, err error) {
	// Вычисляем хэш нового файла
	newFileHash := sha256.Sum256(newFile)
	newFileHashStr := hex.EncodeToString(newFileHash[:])
	newFileSize := int64(len(newFile))

	for _, filePath := range fileNames {
		// Открываем существующий файл
		file, err := os.Open(filePath)
		if err != nil {
			return "", fmt.Errorf("не удалось открыть файл %s: %w", filePath, err)
		}

		// Получаем информацию о файле
		info, err := file.Stat()
		if err != nil {
			file.Close()
			return "", fmt.Errorf("не удалось получить информацию о файле %s: %w", filePath, err)
		}

		// Сравниваем размер файлов
		if info.Size() != newFileSize {
			file.Close()
			continue
		}

		// Читаем содержимое существующего файла
		existingFileContent, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return "", fmt.Errorf("не удалось прочитать файл %s: %w", filePath, err)
		}

		// Вычисляем хэш существующего файла
		existingFileHash := sha256.Sum256(existingFileContent)
		existingFileHashStr := hex.EncodeToString(existingFileHash[:])

		// Сравниваем хэши
		if newFileHashStr == existingFileHashStr {
			// Извлекаем UUID из имени файла
			uuid, err := extractUUID(filePath)
			if err != nil {
				return "", fmt.Errorf("не удалось извлечь UUID из файла %s: %w", filePath, err)
			}
			return uuid, nil
		}
	}

	// Если не найдено совпадений
	return "", nil
}

// extractUUID предполагает, что UUID находится в начале имени файла, разделённого символом '_'
// Например: "123e4567-e89b-12d3-a456-426614174000_filename.ext"
func extractUUID(filePath string) (string, error) {
	// UUID состоит из 36 символов (32 цифры и 4 дефиса)
	runes := []rune(filePath)
	if len(runes) < 36 {
		return "", fmt.Errorf("invalid uuid length")
	}
	return string(runes[:36]), nil
}

func FormFile(r http.Request) (file *models.UploadedFile, err error) {
	r.ParseMultipartForm(10 << 20)

	fileContent, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("FormFile: %w", err)
	}

	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("FormFile (seek 1): %w", err)
	}

	file = &models.UploadedFile{}

	// Переходим в конец файла, чтобы определить его размер
	size, err := fileContent.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("FormFile (seek 2): %w", err)
	}

	// Возвращаемся в исходную позицию чтения
	_, err = fileContent.Seek(currentPos, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("FormFile (seek 3): %w", err)
	}

	// Предварительно аллоцируем слайс байтов нужного размера
	file.Content = make([]byte, size)

	// Читаем содержимое файла в слайс
	_, err = io.ReadFull(fileContent, file.Content)
	if err != nil {
		return nil, fmt.Errorf("FormFile (read): %w", err)
	}

	file.OriginalName = fileHeader.Filename
	file.FileExtension = ExtractFileExtension(file.OriginalName)

	return file, nil

}
