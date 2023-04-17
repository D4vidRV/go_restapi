package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"not null" json:"first_name" validate:"required"`
	LastName  string `gorm:"not null" json:"last_name" validate:"required"`
	Email     string `gorm:"not null;uniqueIndex" json:"email" validate:"required,email"`
	Tasks     []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"tasks"`
}
