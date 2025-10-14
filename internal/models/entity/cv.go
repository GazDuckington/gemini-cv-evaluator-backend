package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CV struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Title     string    `gorm:"size:150;not null" json:"title"`
	FilePath  string    `gorm:"size:255;not null" json:"file_path"`
	Summary   string    `gorm:"type:text" json:"summary"`
	Embedding []float32 `gorm:"type:jsonb" json:"embedding"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (CV) TableName() string {
	return "cvs"
}

func (cv *CV) BeforeCreate(tx *gorm.DB) (err error) {
	cv.ID = uuid.NewString()
	cv.CreatedAt = time.Now()
	return nil
}

func (cv *CV) BeforeUpdate(tx *gorm.DB) (err error) {
	cv.UpdatedAt = time.Now()
	return nil
}
