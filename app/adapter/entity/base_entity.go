package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	Id string `json:"id" gorm:"primaryKey"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	id := uuid.New()
	base.Id = id.String()
	return nil
}
