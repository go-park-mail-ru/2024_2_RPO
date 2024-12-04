package repository

import (
	"RPO_back/internal/models"
	"RPO_back/internal/pkg/utils/uploads"
	"context"
)

func (r *BoardRepository) DeduplicateFile(ctx context.Context, file *models.UploadedFile) (fileNames []string, fileIDs []int64, err error) {
	return uploads.DeduplicateFile(ctx, r.db, file)
}
func (r *BoardRepository) RegisterFile(ctx context.Context, file *models.UploadedFile) (fileID int64, fileUUID string, err error) {
	return uploads.RegisterFile(ctx, r.db, file)
}
