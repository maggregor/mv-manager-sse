package broadcaster

import (
	"log"
	"os"
	"testing"
)

func TestJwtValidationCorrectSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	defer os.Unsetenv("JWT_SECRET")
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.Nv24hvNy238QMrpHvYw-BxyCp00jbsTqjVgzk81PiYA"
	claims, ok := validateAchilioJWT(tokenString)
	if !ok {
		log.Fatalf("Should be true")
	}
	if claims["foo"] != "bar" {
		log.Fatalf("Claim foo should be bar, is %v", claims["foo"])
	}
}

func TestJwtValidationIncorrectSecret(t *testing.T) {
	os.Setenv("JWT_SECRET", "wrongSecret")
	defer os.Unsetenv("JWT_SECRET")
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.Nv24hvNy238QMrpHvYw-BxyCp00jbsTqjVgzk81PiYA"

	_, ok := validateAchilioJWT(tokenString)
	if ok {
		log.Fatalf("Should be false")
	}
}
