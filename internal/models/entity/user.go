package entity

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:100;not null" json:"password_hash"`
	Role         string         `gorm:"size:20;default:user" json:"role"`
	Profile      *Profile       `gorm:"constraint:OnDelete:CASCADE" json:"profile"`
	CVs          []CV           `gorm:"constraint:OnDelete:CASCADE" json:"cvs"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now()
	u.PasswordHash = string(hashed)
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}
