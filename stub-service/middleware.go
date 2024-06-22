package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func RandomErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rand.Seed(time.Now().UnixNano())
		if rand.Intn(100) < 20 {
			c.JSON(500, ErrorResponse{
				Code:    "internal_error",
				Message: "internal server error",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
