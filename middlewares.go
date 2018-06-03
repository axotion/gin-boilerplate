package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func authMiddleware(guard GuardInterface) gin.HandlerFunc {
	return func(c *gin.Context) {

		err := guard.Check(c)

		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			c.Abort()
		}
		c.Next()
	}
}

func throttleMiddleware(limit int) gin.HandlerFunc {
	return func(c *gin.Context) {

		lock.Lock()
		_, ok := blacklisted_ip[c.ClientIP()]

		if !ok {
			blacklisted_ip[c.ClientIP()] = 1
		}
		lock.Unlock()

		blacklisted_ip[c.ClientIP()] = blacklisted_ip[c.ClientIP()] + 1

		if blacklisted_ip[c.ClientIP()] >= 10 {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Calm down. Too many requests!",
			})
			c.Abort()
			return
		}

		go func(list map[string]int, ip string, asyncLock *sync.Mutex) {
			time.Sleep(time.Second * time.Duration(limit))
			asyncLock.Lock()
			list[ip] = 0
			asyncLock.Unlock()

		}(blacklisted_ip, c.ClientIP(), &lock)

		c.Next()

	}
}
