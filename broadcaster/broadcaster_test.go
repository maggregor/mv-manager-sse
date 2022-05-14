package broadcaster

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type configTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}

func (s *configTestSuite) TearDownTest() {
	os.Clearenv()
}

func (s *configTestSuite) TestConfig() {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("SA_EMAIL", "sa-email@gcp.com")
	os.Setenv("AUDIENCE", "http://localhost:8080/events")

	c, err := config()
	s.Nil(err)
	s.Equal("secret", c.JwtSecret)
	s.Equal("sa-email@gcp.com", c.SAEmail)
	s.Equal("http://localhost:8080/events", c.Audience)
}

func (s *configTestSuite) TestConfigError1() {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("SA_EMAIL", "sa-email@gcp.com")

	c, err := config()
	s.NotNil(err)
	s.Equal("secret", c.JwtSecret)
	s.Equal("sa-email@gcp.com", c.SAEmail)
	s.Equal("", c.Audience)
}

func (s *configTestSuite) TestConfigError2() {
	os.Setenv("JWT_SECRET", "")

	c, err := config()
	s.NotNil(err)
	s.Equal("", c.JwtSecret)
}
