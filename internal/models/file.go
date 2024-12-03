package models

type UploadedFile struct {
	Content      []byte // Полное содержимое файла
	OriginalName string // Оригинальное имя файла с расширением
	Hash         string
}
