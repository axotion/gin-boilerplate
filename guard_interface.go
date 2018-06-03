package main

import "github.com/gin-gonic/gin"

type GuardInterface interface {
	SetResolveBy(field string)
	SetSecret(secret []byte)
	SetUserProvider(userProvider UserProviderInterface)
	Attempt(c *gin.Context) error
	Check(c *gin.Context) error
	GetToken() string
	CreateSignedTokenFromID(id int) (string, error)
}
