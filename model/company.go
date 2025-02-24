package model

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name      string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
