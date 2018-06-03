package main

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

type MysqlUserProviderStruct struct {
	db *gorm.DB
}

func (provider *MysqlUserProviderStruct) GetUserById(user int) (*User, error) {
	var tmpUser User
	provider.db.Where("ID = ?", user).First(&tmpUser)

	if tmpUser.ID != 0 {
		return &tmpUser, nil
	}

	return &tmpUser, errors.New("User not found")
}
func (provider *MysqlUserProviderStruct) GetUserBy(username string, password string) (*User, error) {

	var tmpUser User
	var saltUser User

	log.Println(username, password)

	provider.db.Where("email = ?", username).First(&saltUser)

	if saltUser.ID == 0 {
		return nil, errors.New("User not found")
	}

	provider.db.Where("email = ?", username).Where("password = ?", encrypt(saltUser.Salt+password)).First(&tmpUser)

	if tmpUser.ID != 0 {
		return &tmpUser, nil
	}

	return &tmpUser, errors.New("User not found")

}

func (provider *MysqlUserProviderStruct) IsEmailExist(email string) bool {
	var tmpUser User
	provider.db.Where("email = ?", email).First(&tmpUser)

	if tmpUser.ID != 0 {
		return true
	}

	return false
}
