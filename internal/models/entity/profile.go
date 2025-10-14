package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	FullName  string `gorm:"size:100;not null" json:"full_name"`
	Bio       string `gorm:"type:text" json:"bio"`
	Phone     string `gorm:"size:20" json:"phone"`
	AvatarURL string `gorm:"size:255" json:"avatar_url"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Profile) TableName() string {
	return "profiles"
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	p.CreatedAt = time.Now()
	return nil
}

func (p *Profile) BeforeUpdate(tx *gorm.DB) (err error) {
	p.UpdatedAt = time.Now()
	return nil
}
