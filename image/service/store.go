package image_service

import (
	"context"
	"heyalley-server/db"
	"mime/multipart"
	"os"
	"path/filepath"
)

type ImageStore interface {
	StoreImage(context.Context, multipart.File, *multipart.FileHeader) error
	GetImage(ctx context.Context, searchKey string) (string, error)
}

type ImageStoreService struct {
	DatabaseService       db.DatabaseService
	ImageRegisteryService ImageRegistery
}

func (s *ImageStoreService) StoreImage(ctx context.Context, file multipart.File, h *multipart.FileHeader) error {
	newFileName, err := s.ImageRegisteryService.RegisterImageInDatabase(ctx, h.Filename)
	if err != nil {
		return err
	}
	err = s.ImageRegisteryService.RegisterImageInFileSystem(ctx, file, newFileName)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImageStoreService) GetImage(ctx context.Context, searchKey string) (string, error) {

	currDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}
	imagesDirectory := filepath.Join(currDirectory, "images-blob")

	image, err := s.DatabaseService.GetBySearchKey(ctx, searchKey)
	if err != nil {
		return "", err
	}
	return imagesDirectory + "/" + image.Path, nil
}
