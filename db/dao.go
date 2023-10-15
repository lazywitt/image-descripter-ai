package db

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"

	"heyalley-server/db/models"
)

type Handler struct {
	DB *gorm.DB
}

func GetHandler(dbClient *gorm.DB) *Handler {
	return &Handler{
		DB: dbClient,
	}
}

type DatabaseService interface {
	// Create inserts the given images into entity
	Create(ctx context.Context, images []*models.Image) error
	// GetById fetches entry by Id
	GetById(ctx context.Context, id string) (*models.Image, error)
	// GetByToken fetches paginated response wrt given token
	// GetByToken(ctx context.Context, token string) (*models.Images, error)
	// GetBySearchKey returns images which has description matching with given searchKey
	// matching results is performed by matching every token in searchKey with description
	// Example: "new hat" will match "hat in new york" and "old hat and new hat" both. match token are being
	// created by description. Query size limited by 10
	GetBySearchKey(ctx context.Context, searchKey string) (*models.Image, error)
}

func (h *Handler) Create(ctx context.Context, images []*models.Image) error {

	res := h.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(images)
	if res.Error != nil {
		return fmt.Errorf("error creating entries in DB %v", res.Error)
	}
	return nil
}

func (h *Handler) UpdateImage(image *models.Image) error {
	res := h.DB.Updates(image)
	if res.Error != nil {
		return fmt.Errorf("error updating image data: %v", res.Error)
	}
	return nil
}

func (h *Handler) GetUnprocessedImage(ctx context.Context) (*models.Image, error) {
	var (
		res = &models.Image{}
	)
	gormRes := h.DB.Where("is_pipeline_processed = false").Order("created_at").Limit(1).First(res)
	if gormRes.Error != nil {
		return nil, gormRes.Error
	}
	return res, gormRes.Error
}

func (h *Handler) GetById(ctx context.Context, id string) (*models.Image, error) {
	var (
		getByTokenRes = &models.Image{}
	)
	gormRes := h.DB.Where("Id = ?", id).First(getByTokenRes)
	return getByTokenRes, gormRes.Error
}

func (h *Handler) GetBySearchKey(ctx context.Context, searchKey string) (*models.Image, error) {
	var (
		modelimage *models.Image = &models.Image{}
	)

	gormRes := h.DB.Raw("select * from images where to_tsvector(description) @@ to_tsquery( ? ) limit 1;",
		strings.ReplaceAll(searchKey, " ", "&")).First(modelimage)

	if gormRes.Error != nil {
		return nil, gormRes.Error
	}
	return modelimage, nil
}
