package main

import (
	"net/http"

	"tgs-automation/internal/log"

	jwttoken "tgs-automation/internal/jwt_token"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// authMiddleware validates the JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BearerSchema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			return
		}

		tokenString := header[len(BearerSchema):]
		log.LogInfo("Token: " + tokenString)
		token, err := jwt.ParseWithClaims(tokenString, &jwttoken.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwttoken.JwtKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Next()
	}
}
