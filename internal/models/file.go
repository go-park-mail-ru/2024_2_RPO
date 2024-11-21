package models

type UploadedFile struct {
	FileID        *int64 // Суррогатный первичный ключ в таблице user_uploaded_file
	Content       []byte // Полное содержимое файла
	OriginalName  string
	UUID          *string // UUID файла (в таблице user_uploaded_file)
	FileExtension string  // nil, если у файла нет расширения
}
