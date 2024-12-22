package models

//go:generate easyjson -all file.go

type UploadedFile struct {
	Content      []byte // Полное содержимое файла
	OriginalName string // Оригинальное имя файла с расширением
	Hash         string
}
