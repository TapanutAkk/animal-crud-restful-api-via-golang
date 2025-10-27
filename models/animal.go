package models

import "gorm.io/gorm"

type Animal struct {
	gorm.Model
	Name      string `json:"name" gorm:"not null"`
	Species   string `json:"species"`
	Age       int    `json:"age"`
}