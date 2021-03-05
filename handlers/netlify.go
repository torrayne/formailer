package handlers

import (
	"io"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/djatwood/formailer"
)

func netlifyResponse(code int, err error, headers ...[2]string) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       http.StatusText(code),
		Headers:    make(map[string]string),
	}

	for _, h := range headers {
		response.Headers[h[0]] = h[1]
	}

	if err != nil {
		response.Body = err.Error()
		if _, ok := response.Headers["location"]; ok {
			response.Headers["location"] += "?error=" + err.Error()
		}
	}

	return response
}

// Netlify takes in a aws lambda request and sends an email
func Netlify(c *formailer.Config) func(events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		if request.HTTPMethod != "POST" {
			return netlifyResponse(http.StatusMethodNotAllowed, nil), nil
		}

		submission, err := c.Parse(request.Headers["content-type"], request.Body)
		if err != nil && err != io.EOF {
			return netlifyResponse(http.StatusBadRequest, err), nil
		}

		if v, ok := submission.Values["g-recaptcha-response"]; ok {
			ok, err := verifyRecaptcha(v.(string))
			if err != nil {
				return netlifyResponse(http.StatusInternalServerError, err), nil
			}
			if !ok {
				return netlifyResponse(http.StatusBadRequest, nil), nil
			}

			delete(submission.Values, "g-recaptcha-response")
		}

		server, err := submission.Form.SMTPServer()
		if err != nil {
			return netlifyResponse(500, nil), nil
		}

		err = submission.Send(server)
		if err != nil {
			return netlifyResponse(http.StatusInternalServerError, err), nil
		}

		statusCode := http.StatusOK
		headers := [][2]string{}
		if len(submission.Form.Redirect) > 0 {
			statusCode = http.StatusSeeOther
			headers = append(headers, [2]string{"location", submission.Form.Redirect})
		}

		return netlifyResponse(statusCode, nil, headers...), nil
	}
}
