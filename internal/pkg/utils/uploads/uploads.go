package uploads

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/config"
	"RPO_back/internal/pkg/utils/logging"
	"RPO_back/internal/pkg/utils/pgxiface"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
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

// CompareFiles смотрит, равен ли данный файл какому-нибудь из существующих
// загруженных файлов. Если да, возвращает fileUUID эквивалентного файла
func CompareFiles(fileNames []string, fileIDs []int64, newFile *models.UploadedFile) (fileID *int64, err error) {
	for idx, filePath := range fileNames {
		// Читаем существующий файл
		file, err := os.Open(filepath.Join(config.CurrentConfig.UploadsDir, filePath))
		if err != nil {
			return nil, fmt.Errorf("CompareFiles (open) %s: %w", filePath, err)
		}
		existingFileContent, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("CompareFiles (read) %s: %w", filePath, err)
		}

		// Сравнение по содержимому
		if slices.Equal(newFile.Content, existingFileContent) {
			return &fileIDs[idx], err
		}
	}

	// Не найдено совпадений
	return nil, nil
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

func FormFile(r *http.Request) (file *models.UploadedFile, err error) {
	r.ParseMultipartForm(10 << 20)

	fileContent, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("FormFile: %w", err)
	}

	currentPos, err := fileContent.Seek(0, io.SeekCurrent)
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

func SaveFile(file *models.UploadedFile) (err error) {
	if file.UUID == nil {
		return fmt.Errorf("SaveFile: file uuid is nil")
	}
	fileName := JoinFilePath(*file.UUID, file.FileExtension)
	err = os.WriteFile(filepath.Join(config.CurrentConfig.UploadsDir, fileName), file.Content, 0644)
	if err != nil {
		return fmt.Errorf("SaveFile: cant write to file %s: %w", fileName, err)
	}
	return nil
}

// DeduplicateFile возвращает список файлов с таким же расширением, ID и размером
func DeduplicateFile(ctx context.Context, db pgxiface.PgxIface, file *models.UploadedFile) (fileNames []string, fileIDs []int64, err error) {
	funcName := "DeduplicateFile"
	query := `
	SELECT file_id, file_uuid, file_extension
	FROM user_uploaded_file
	WHERE file_hash=$1 AND "size"=$2 AND file_extension=$3;
	`

	h := sha1.New()
	h.Write(file.Content)
	fileHash := hex.EncodeToString(h.Sum(nil))

	rows, err := db.Query(ctx, query, fileHash, len(file.Content), file.FileExtension)
	if err != nil {
		return nil, nil, fmt.Errorf("%s (query): %w", funcName, err)
	}
	for rows.Next() {
		var fileUUID, fileExtension string
		var fileID int64
		err = rows.Scan(&fileID, &fileUUID, &fileExtension)
		if err != nil {
			return nil, nil, fmt.Errorf("%s (scan): %w", funcName, err)
		}
		fileIDs = append(fileIDs, fileID)
		fileNames = append(fileNames, JoinFilePath(fileUUID, fileExtension))
	}
	return fileNames, fileIDs, nil
}

// RegisterFile заносит информацию о файле в таблицу и по указателю меняет поля FileID и UUID в структуре file
func RegisterFile(ctx context.Context, db pgxiface.PgxIface, file *models.UploadedFile) error {
	funcName := "RegisterFile"
	query := `
	INSERT INTO user_uploaded_file
	(file_extension, created_at, "size")
	VALUES ($1, CURRENT_TIMESTAMP, $2)
	RETURNING file_uuid::text, file_id;
	`
	row := db.QueryRow(ctx, query, file.FileExtension, len(file.Content))
	err := row.Scan(&file.UUID, file.FileID)
	logging.Debug(ctx, funcName, " query has err: ", err)
	if err != nil {
		return fmt.Errorf("%s: %w", funcName, err)
	}
	return nil
}
