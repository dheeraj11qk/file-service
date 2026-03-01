package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BeforeCreateInterface interface {
	BeforeCreate(*gorm.DB) error
}

type File struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	FileName   string
	FileSize   int64
	FilePath   string
	UploadTime time.Time
	Status     string
}

func (f *File) BeforeCreate(tx *gorm.DB) error {
	f.ID = uuid.New()
	return nil
}
