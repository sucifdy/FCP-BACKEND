package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var secretKey = model.JwtKey

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil token dari cookie dengan nama "session_token"
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			if ctx.Request.URL.Path != "/" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				ctx.Abort()
				return
			} else {
				ctx.Redirect(http.StatusSeeOther, "/login")
				ctx.Abort()
				return
			}
		}

		tokenString := cookie.Value
		log.Println("Token received:", tokenString)

		// Parse token dengan claims
		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if err != nil {
			log.Println("Error parsing token:", err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}

		// Periksa jika token valid dan set email di context
		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
			log.Println("Token valid, Email:", claims.Email)
			ctx.Set("email", claims.Email)
			ctx.Next()
		} else {
			log.Println("Invalid token")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		}
	}
}
