package broadcaster

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("SA_EMAIL", "sa-email@gcp.com")
	os.Setenv("AUDIENCE", "http://localhost/events")

	c := config()
	assert.Equal(t, "secret", c.JwtSecret)
	assert.Equal(t, "sa-email@gcp.com", c.SAEmail)
	assert.Equal(t, "http://localhost/events", c.Audience)
}
