package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Image model entity
type Image struct {
	Id                  uuid.UUID      `json:"id" gorm:"primaryKey"`
	Path                string         `json:"youtubeId" gorm:"unique"`
	OriginalFileName    sql.NullString `json:"originalFileName"`
	Description         sql.NullString `json:"description"`
	IsPipelineProcessed bool           `json:"isPipelineProcessed" gorm:"default:false"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           time.Time      `json:"deletedAt" gorm:"default:null"`
}
