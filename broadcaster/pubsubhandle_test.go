package broadcaster

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type pubSubTestSuite struct {
	suite.Suite
}

func TestPubSubTestSuite(t *testing.T) {
	suite.Run(t, new(pubSubTestSuite))
}

func (s *pubSubTestSuite) SetupSuite() {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("ENV", "local")
	os.Setenv("PORT", "8181")
	go Serve()
	time.Sleep(1 * time.Second)
}

func (s *pubSubTestSuite) TestHandPubSub() {
	body := strings.NewReader(`{"deliveryAttempt":1,"message":{"data":"eyJzdGF0dXMiOiJDT01QTEVURUQifQ==","attributes":{"eventType":"NEW_QUERIES","projectId":"achilio-dev","teamName":"achilio.com"},"messageId":"4432805535655886","message_id":"4432805535655886","publishTime":"2022-04-28T15:57:52.096Z","publish_time":"2022-04-28T15:57:52.096Z"},"subscription":"projects/achilio-dev/subscriptions/mvSseSubscription"}`)
	req, err := http.NewRequest("POST", "http://localhost:8181/events", body)
	s.NoError(err)
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
}
