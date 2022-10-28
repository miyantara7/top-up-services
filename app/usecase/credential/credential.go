package credential

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const secret = "secret"

func GenerateToken(username string) (t string, e error) {

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return t, nil
}
