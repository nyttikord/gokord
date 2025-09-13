package gokord

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/nyttikord/gokord/discord"
)

// All error constants
var (
	ErrJSONUnmarshal = errors.New("json unmarshal")
	ErrStatusOffline = errors.New("you can't set your Status to offline")
	ErrUnauthorized  = errors.New("HTTP request was unauthorized. This could be because the provided token was not a bot token. Please add \"Bot \" to the start of your token. https://discord.com/developers/docs/reference#authentication-example-bot-token-authorization-header")
)

// RESTError stores error information about a request with a bad response code.
// Message is not always present, there are cases where api calls can fail
// without returning a json message.
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte

	// Message may be nil.
	Message *discord.APIErrorMessage
}

// newRestError returns a new REST API error.
func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	// Attempt to decode the error and assume no message was provided if it fails
	var msg *discord.APIErrorMessage
	err := json.Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

// Error returns a Rest API Error with its status code and body.
func (r RESTError) Error() string {
	base := fmt.Sprintf("[HTTP %d]", r.Response.StatusCode)
	if r.Message != nil {
		return fmt.Sprintf("%s %s", base, r.Message.Error())
	}
	return fmt.Sprintf("%s %s", base, r.ResponseBody)
}

// RateLimitError is returned when a request exceeds a rate limit and Session.ShouldRetryOnRateLimit is false.
// The request may be manually retried after waiting the duration specified by RetryAfter.
type RateLimitError struct {
	*RateLimit
}

// Error returns a rate limit error with rate limited endpoint and retry time.
func (e RateLimitError) Error() string {
	return "Rate limit exceeded on " + e.URL + ", retrying after " + e.RetryAfter.String()
}
