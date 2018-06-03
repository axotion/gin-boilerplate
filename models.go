package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string     `json:"name" binding:"required"`
	Birthday *time.Time `json:"birthday" binding:"required"`
	Email    string     `gorm:"type:varchar(100);unique_index" json:"email" binding:"required"`
	Type     string     `gorm:"size:255" json:"type" binding:"required,allowed_type"`
	Password string     `gorm:"size:255" json:"password,omitempty" binding:"required"`
	Salt     string     `gorm:"size:255" json:"-"`
	Address  string     `gorm:"index:addr" json:"address" binding:"required"`
	Profile  Profile    `gorm:"foreignkey:UserID;association_foreignkey:Refer" json:"-"`
}

type Profile struct {
	gorm.Model
	Links  []string `gorm:"type:varchar(100);" json:"-"`
	UserID uint
	Refer  int
}

type LoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
