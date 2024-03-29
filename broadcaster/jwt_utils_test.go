package broadcaster

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwtValidationCorrectSecret(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.Nv24hvNy238QMrpHvYw-BxyCp00jbsTqjVgzk81PiYA"

	claims, err := validateAchilioJWT(tokenString, "secret")
	assert.Nil(t, err)
	assert.Equal(t, "bar", claims["foo"])
}

func TestJwtValidationIncorrectSecret(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.Nv24hvNy238QMrpHvYw-BxyCp00jbsTqjVgzk81PiYA"

	_, err := validateAchilioJWT(tokenString, "wrongSecret")
	assert.NotNil(t, err, "Should be in error")
}

func TestJwtValidationSignatureAlg(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.Nv24hvNy238QMrpHvYw-BxyCp00jbsTqjVgzk81PiYA"

	_, err := validateAchilioJWT(tokenString, "secret")
	assert.NotNil(t, err, "Should be in error")
}

func TestGetJwtFromCookie(t *testing.T) {
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg"
	cookie1 := "jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg"
	cookie2 := "phpsseid=298zf09hf012fh2;jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg"
	cookie3 := "phpsseid=298zf09hf012fh2; jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg"
	cookie4 := "phpsseid=298zf09hf012fh2; \tjwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg"
	cookie5 := ""
	cookie6 := "phpsseid=298zf09hf012fh2"
	cookie7 := "phpsseid=298zf09hf012fh2;"
	cookie8 := "phpsseid=298zf09hf012fh2,jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg"
	cookie9 := "phpsseid:298zf09hf012fh2"
	assert.Equal(t, expected, getJwtFromCookie(strings.Split(cookie1, ";")))
	assert.Equal(t, expected, getJwtFromCookie(strings.Split(cookie2, ";")))
	assert.Equal(t, expected, getJwtFromCookie(strings.Split(cookie3, ";")))
	assert.Equal(t, expected, getJwtFromCookie(strings.Split(cookie4, ";")))
	assert.Equal(t, "", getJwtFromCookie(strings.Split(cookie5, ";")))
	assert.Equal(t, "", getJwtFromCookie(strings.Split(cookie6, ";")))
	assert.Equal(t, "", getJwtFromCookie(strings.Split(cookie7, ";")))
	assert.Equal(t, "", getJwtFromCookie(strings.Split(cookie8, ";")))
	assert.Equal(t, "", getJwtFromCookie(strings.Split(cookie9, ";")))
}
