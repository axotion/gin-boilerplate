package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type JwtGuard struct {
	secret       []byte
	userProvider UserProviderInterface
	resolveBy    string
	token        string
}

func (j *JwtGuard) SetResolveBy(field string) {
	j.resolveBy = field
}

func (j *JwtGuard) SetSecret(secret []byte) {
	j.secret = secret
}

func (j *JwtGuard) SetUserProvider(userProvider UserProviderInterface) {
	j.userProvider = userProvider
}

func (j *JwtGuard) Attempt(c *gin.Context) error {

	raw_username, ok_username := c.Get(j.resolveBy)
	raw_password, ok_password := c.Get("password")

	if !ok_password || !ok_username {
		return errors.New("invalid parameters")
	}

	username, ok_username := raw_username.(string)
	password, ok_password := raw_password.(string)

	if !ok_password || !ok_username {
		return errors.New("invalid parameters")
	}

	user, err := j.userProvider.GetUserBy(username, password)

	if err != nil {
		return errors.New("Username or password not found")
	}

	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		Issuer:    fmt.Sprint(user.ID),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)

	if err != nil {
		log.Println(err)
		return errors.New("something bad happen")
	}

	j.token = signedToken

	return nil
}

func (j *JwtGuard) Check(c *gin.Context) error {

	stringToken := c.GetHeader("Authorization")
	j.token = stringToken

	if stringToken == "" {
		return errors.New("No auth header provided")
	}

	authHeader := strings.Split(stringToken, " ")

	if len(authHeader) != 2 {
		return errors.New("No auth header provided")
	}

	stringToken = authHeader[1]

	if stringToken == "" {
		return errors.New("No auth header provided")
	}

	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.secret), nil
	})

	if err != nil {
		return errors.New(err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {

		UserID, err := strconv.Atoi(claims["iss"].(string))

		if err != nil {
			return errors.New("Invalid token")
		}

		log.Println(UserID)

		user, err := j.userProvider.GetUserById(UserID)

		if err != nil {
			return err
		}
		c.Set("user", user)
		return nil
	}
	return err
}

func (j *JwtGuard) GetToken() string {
	return j.token
}

func (j *JwtGuard) CreateSignedTokenFromID(id int) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		Issuer:    fmt.Sprint(id),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(j.secret)

	if err != nil {
		return "", errors.New(err.Error())
	}

	return signedToken, nil
}
