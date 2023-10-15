package image_service

import (
	"context"
	"database/sql"
	"heyalley-server/db"
	"heyalley-server/db/models"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/google/uuid"
)

type ImageRegistery interface {
	RegisterImageInFileSystem(context.Context, multipart.File, string) error
	RegisterImageInDatabase(context.Context, string) (string, error)
}

type ImageRegisteryService struct {
	DatabaseService db.DatabaseService
}

func (s *ImageRegisteryService) RegisterImageInFileSystem(ctx context.Context, file multipart.File, fileName string) error {
	path := "./images-blob/" + fileName
	tmpfile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer tmpfile.Close()

	_, err = io.Copy(tmpfile, file)
	if err != nil {
		return err
	}
	return nil
}

func (s *ImageRegisteryService) RegisterImageInDatabase(ctx context.Context, fileName string) (string, error) {
	newFileName := time.Now().Format("20060102150405") + fileName
	nm := &models.Image{
		Id:   uuid.New(),
		Path: newFileName,
	}
	if fileName != "" {
		nm.OriginalFileName = sql.NullString{String: fileName, Valid: true}
	}
	err := s.DatabaseService.Create(ctx, []*models.Image{nm})
	if err != nil {
		return "", err
	}
	return newFileName, nil
}

// Image is api level response structure
// type Image struct {
// 	Id          string
// 	Path        string
// 	FileName    string
// 	Description string
// }

// type GetVideoResponse struct {
// 	Videos    []Video `json:"videos"`
// 	NextToken string  `json:"token"`
// }onse struct {
// 	Videos    []Video `json:"vi
