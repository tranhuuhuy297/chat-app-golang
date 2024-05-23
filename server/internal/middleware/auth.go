package middleware

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	log "github.com/sirupsen/logrus"
)

var SECRET_KEY string

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	SECRET_KEY = os.Getenv("SECRET_KEY")
}

func AuthenticateMiddleware(c *gin.Context) {
	tokenString := c.GetHeader("Bearer")
	log.Printf("middleware/auth| token: %s\n", tokenString)

	token, err := verifyToken(tokenString)
	if err != nil {
		log.Printf("middleware/auth| token error: %v\n", err)
	}
	log.Printf("middleware/auth| parse token: %v\n", token.Claims)
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//TODO check HS256 or RS256 jwt signing, default is HS256
		return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		log.Printf("middleware/auth| check token.Valid \n")
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}
