package main

type UserProviderInterface interface {
	GetUserById(user int) (*User, error)
	GetUserBy(username string, password string) (*User, error)
	IsEmailExist(email string) bool
}
