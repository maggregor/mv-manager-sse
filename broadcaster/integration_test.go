package broadcaster

import (
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type e2eTestSuite struct {
	suite.Suite
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(e2eTestSuite))
}

func (s *e2eTestSuite) SetupSuite() {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("ENV", "paslocal")
	os.Setenv("SA_EMAIL", "foo")
	os.Setenv("AUDIENCE", "bar")
	os.Setenv("PORT", "8180")
	go Serve()
	time.Sleep(1 * time.Second)
}

func (s *e2eTestSuite) TestSubscribe() {
	req, err := http.NewRequest("GET", "http://localhost:8180/subscribe", nil)
	s.NoError(err)
	req.Header.Set("Cookie", "jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkCg")
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	s.Equal("text/event-stream", response.Header.Get("content-type"))
}

func (s *e2eTestSuite) TestSubscribe_InvalidJWT() {
	req, err := http.NewRequest("GET", "http://localhost:8180/subscribe", nil)
	s.NoError(err)
	req.Header.Set("Cookie", "jwt_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJoZCI6ImFjaGlsaW8uY29tIiwibmJmIjoxNDQ0NDc4NDAwfQ.2BrdTRK_qQmcGiXe-ApEey0-fm0xvBK90ssjFmvHkC")
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusUnauthorized, response.StatusCode)

	req.Header.Set("Cookie", "csrftoken=abcdef")
	response, err = client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusUnauthorized, response.StatusCode)
}

func (s *e2eTestSuite) TestValidatePubSubJwt_InvalidJWT() {
	body := strings.NewReader(`{"deliveryAttempt":1,"message":{"attributes":{"eventType":"NEW_QUERIES","projectId":"achilio-dev","teamName":"achilio.com"},"messageId":"4432805535655886","message_id":"4432805535655886","publishTime":"2022-04-28T15:57:52.096Z","publish_time":"2022-04-28T15:57:52.096Z"},"subscription":"projects/achilio-dev/subscriptions/mvSseSubscription"}`)
	req, err := http.NewRequest("POST", "http://localhost:8180/events", body)
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg2MTY0OWU0NTAzMTUzODNmNmI5ZDUxMGI3Y2Q0ZTkyMjZjM2NkODgiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwczovL2Rldi5zLmFjaGlsaW8uY29tL2V2ZW50cyIsImF6cCI6IjEwODc0ODMyMDQ2MjIyMjQ4MTE3NyIsImVtYWlsIjoic3NlLWNsb3VkcnVuQGFjaGlsaW8tZGV2LmlhbS5nc2VydmljZWFjY291bnQuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MTY1MTE2NDMxOSwiaWF0IjoxNjUxMTYwNzE5LCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiIxMDg3NDgzMjA0NjIyMjI0ODExNzcifQ.BH9As14ucCjullb3muOFWR2RVv6Q4SPzRH_Cx2WCtIkAh7fefAsv4qtUuj6X6YRmCpITxYDVz6WenAz6tDS-5GkC6399vU5AQartwAUVy_AVj1g6LdeEBJ-0ND6tEfTLTnHIVNunLmyMyi01ky9D3p_5i6qCHtA2AZN8jvbRLMxIRKgePyErU0FPbBJ76dAFQvj35-yORgCCL-sado3K_PyOD9Ll8PpxiUIuGAv7nOxJ1D0sq1DrxCIax1794UpR8RZf7qhs1gONfdvgobbrceRFasYI4D19KNYsPVTWMAF8vtKRjG0SorldMUAnsGe2luY3NH16t16Qe8jGTxgguw")
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusBadRequest, response.StatusCode)
}

func (s *e2eTestSuite) TestValidatePubSubJwt_InvalidMethod() {
	body := strings.NewReader(`{"deliveryAttempt":1,"message":{"attributes":{"eventType":"NEW_QUERIES","projectId":"achilio-dev","teamName":"achilio.com"},"messageId":"4432805535655886","message_id":"4432805535655886","publishTime":"2022-04-28T15:57:52.096Z","publish_time":"2022-04-28T15:57:52.096Z"},"subscription":"projects/achilio-dev/subscriptions/mvSseSubscription"}`)
	req, err := http.NewRequest("GET", "http://localhost:8180/events", body)
	s.NoError(err)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg2MTY0OWU0NTAzMTUzODNmNmI5ZDUxMGI3Y2Q0ZTkyMjZjM2NkODgiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwczovL2Rldi5zLmFjaGlsaW8uY29tL2V2ZW50cyIsImF6cCI6IjEwODc0ODMyMDQ2MjIyMjQ4MTE3NyIsImVtYWlsIjoic3NlLWNsb3VkcnVuQGFjaGlsaW8tZGV2LmlhbS5nc2VydmljZWFjY291bnQuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MTY1MTE2NDMxOSwiaWF0IjoxNjUxMTYwNzE5LCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiIxMDg3NDgzMjA0NjIyMjI0ODExNzcifQ.BH9As14ucCjullb3muOFWR2RVv6Q4SPzRH_Cx2WCtIkAh7fefAsv4qtUuj6X6YRmCpITxYDVz6WenAz6tDS-5GkC6399vU5AQartwAUVy_AVj1g6LdeEBJ-0ND6tEfTLTnHIVNunLmyMyi01ky9D3p_5i6qCHtA2AZN8jvbRLMxIRKgePyErU0FPbBJ76dAFQvj35-yORgCCL-sado3K_PyOD9Ll8PpxiUIuGAv7nOxJ1D0sq1DrxCIax1794UpR8RZf7qhs1gONfdvgobbrceRFasYI4D19KNYsPVTWMAF8vtKRjG0SorldMUAnsGe2luY3NH16t16Qe8jGTxgguw")
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusMethodNotAllowed, response.StatusCode)
}

func (s *e2eTestSuite) TestValidatePubSubJwt_InvalidHeader1() {
	body := strings.NewReader(`{"deliveryAttempt":1,"message":{"attributes":{"eventType":"NEW_QUERIES","projectId":"achilio-dev","teamName":"achilio.com"},"messageId":"4432805535655886","message_id":"4432805535655886","publishTime":"2022-04-28T15:57:52.096Z","publish_time":"2022-04-28T15:57:52.096Z"},"subscription":"projects/achilio-dev/subscriptions/mvSseSubscription"}`)
	req, err := http.NewRequest("POST", "http://localhost:8180/events", body)
	s.NoError(err)
	req.Header.Set("Authorizations", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg2MTY0OWU0NTAzMTUzODNmNmI5ZDUxMGI3Y2Q0ZTkyMjZjM2NkODgiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwczovL2Rldi5zLmFjaGlsaW8uY29tL2V2ZW50cyIsImF6cCI6IjEwODc0ODMyMDQ2MjIyMjQ4MTE3NyIsImVtYWlsIjoic3NlLWNsb3VkcnVuQGFjaGlsaW8tZGV2LmlhbS5nc2VydmljZWFjY291bnQuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MTY1MTE2NDMxOSwiaWF0IjoxNjUxMTYwNzE5LCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiIxMDg3NDgzMjA0NjIyMjI0ODExNzcifQ.BH9As14ucCjullb3muOFWR2RVv6Q4SPzRH_Cx2WCtIkAh7fefAsv4qtUuj6X6YRmCpITxYDVz6WenAz6tDS-5GkC6399vU5AQartwAUVy_AVj1g6LdeEBJ-0ND6tEfTLTnHIVNunLmyMyi01ky9D3p_5i6qCHtA2AZN8jvbRLMxIRKgePyErU0FPbBJ76dAFQvj35-yORgCCL-sado3K_PyOD9Ll8PpxiUIuGAv7nOxJ1D0sq1DrxCIax1794UpR8RZf7qhs1gONfdvgobbrceRFasYI4D19KNYsPVTWMAF8vtKRjG0SorldMUAnsGe2luY3NH16t16Qe8jGTxgguw")
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusBadRequest, response.StatusCode)
}

func (s *e2eTestSuite) TestValidatePubSubJwt_InvalidHeader2() {
	body := strings.NewReader(`{"deliveryAttempt":1,"message":{"attributes":{"eventType":"NEW_QUERIES","projectId":"achilio-dev","teamName":"achilio.com"},"messageId":"4432805535655886","message_id":"4432805535655886","publishTime":"2022-04-28T15:57:52.096Z","publish_time":"2022-04-28T15:57:52.096Z"},"subscription":"projects/achilio-dev/subscriptions/mvSseSubscription"}`)
	req, err := http.NewRequest("POST", "http://localhost:8180/events", body)
	s.NoError(err)
	req.Header.Set("Authorization", "eyJhbGciOiJSUzI1NiIsImtpZCI6Ijg2MTY0OWU0NTAzMTUzODNmNmI5ZDUxMGI3Y2Q0ZTkyMjZjM2NkODgiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJodHRwczovL2Rldi5zLmFjaGlsaW8uY29tL2V2ZW50cyIsImF6cCI6IjEwODc0ODMyMDQ2MjIyMjQ4MTE3NyIsImVtYWlsIjoic3NlLWNsb3VkcnVuQGFjaGlsaW8tZGV2LmlhbS5nc2VydmljZWFjY291bnQuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImV4cCI6MTY1MTE2NDMxOSwiaWF0IjoxNjUxMTYwNzE5LCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiIxMDg3NDgzMjA0NjIyMjI0ODExNzcifQ.BH9As14ucCjullb3muOFWR2RVv6Q4SPzRH_Cx2WCtIkAh7fefAsv4qtUuj6X6YRmCpITxYDVz6WenAz6tDS-5GkC6399vU5AQartwAUVy_AVj1g6LdeEBJ-0ND6tEfTLTnHIVNunLmyMyi01ky9D3p_5i6qCHtA2AZN8jvbRLMxIRKgePyErU0FPbBJ76dAFQvj35-yORgCCL-sado3K_PyOD9Ll8PpxiUIuGAv7nOxJ1D0sq1DrxCIax1794UpR8RZf7qhs1gONfdvgobbrceRFasYI4D19KNYsPVTWMAF8vtKRjG0SorldMUAnsGe2luY3NH16t16Qe8jGTxgguw")
	client := http.DefaultClient
	response, err := client.Do(req)
	s.NoError(err)
	s.Equal(http.StatusBadRequest, response.StatusCode)
}
