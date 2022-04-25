package broadcaster

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret []byte

func validateAchilioJWT(tokenString string) (jwt.MapClaims, bool) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		hmacSampleSecret = []byte(os.Getenv("JWT_SECRET"))
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Println(err)
		return nil, false
	}
}

func getJwtFromCookie(cookie []string) string {
	for _, c := range cookie {
		slice := strings.Split(c, "=")
		if strings.Trim(slice[0], " \t") == "jwt_token" {
			return slice[1]
		}
	}
	return ""
}
