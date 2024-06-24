package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
)

func RandomErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
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
