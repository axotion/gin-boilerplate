package main

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

func indexAction(c *gin.Context) {
	c.JSON(200, gin.H{
		"API_VERSION": apiVersion,
		"DOC_REFER":   nil,
		"AUTH_METHOD": authMethod,
	})
}

func registerAction(c *gin.Context) {

	var user User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	userProvider := Container.Resolve("user_provider").(UserProviderInterface)
	guard := Container.Resolve("guard").(GuardInterface)
	db := Container.Resolve("database").(*gorm.DB)

	if userProvider.IsEmailExist(user.Email) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Email already taken",
		})
		return
	}

	user.Salt = strconv.Itoa(rand.Int())

	//TODO: Before observer from GORM to encrypt password
	user.Password = encrypt(user.Salt + user.Password)

	db.Create(&user)

	//We dont want to return hashed password
	user.Password = ""

	token, err := guard.CreateSignedTokenFromID(int(user.ID))

	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

func loginAction(c *gin.Context) {

	var loginForm LoginForm

	if err := c.ShouldBind(&loginForm); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	userProvider := Container.Resolve("user_provider").(UserProviderInterface)
	guard := Container.Resolve("guard").(GuardInterface)

	user, err := userProvider.GetUserBy(loginForm.Email, loginForm.Password)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := guard.CreateSignedTokenFromID(int(user.ID))

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"error": err.Error(),
		})
		return
	}

	//We dont want to return
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})

}
