package formailer

//go:generate go run generate/main.go

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aymerick/douceur/inliner"
	mail "github.com/xhit/go-simple-mail/v2"
)

type form struct {
	name    string
	to      string
	from    string
	subject string
}

func setup() (*mail.SMTPServer, error) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Println("could not parse SMTP_PORT")
		return nil, err
	}

	server := mail.NewSMTPClient()
	server.Host = os.Getenv("SMTP_HOST")
	server.Port = port
	server.Username = os.Getenv("SMTP_USER")
	server.Password = os.Getenv("SMTP_PASS")
	server.Encryption = mail.EncryptionTLS
	server.Authentication = mail.AuthLogin
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	return server, nil
}

func respond(code int, err error) *events.APIGatewayProxyResponse {
	response := &events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       http.StatusText(code),
	}

	if err == nil {
		return response
	}

	str, err := json.Marshal(map[string]string{"message": err.Error()})
	if err != nil {
		code = http.StatusInternalServerError
		response.Body = http.StatusText(http.StatusInternalServerError)
	} else {
		response.Body = string(str)
	}

	return response
}

func parseData(contentType, body string) (map[string]string, error) {
	data := make(map[string]string)

	if strings.Contains(contentType, "application/json") {
		err := json.Unmarshal([]byte(body), &data)
		if err != nil {
			return data, err
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		vals, err := url.ParseQuery(body)
		if err != nil {
			return data, err
		}
		for k := range vals {
			data[k] = vals.Get(k)
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		fmt.Println("Content Type", contentType)

		reader := multipart.NewReader(strings.NewReader(body), "formailer")
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			value := new(bytes.Buffer)
			value.ReadFrom(part)
			data[part.FormName()] = value.String()
		}

		fmt.Println(data)
		return data, errors.New("temp")
	} else {
		fmt.Println(contentType)
		return data, errors.New("invalid content type")
	}

	return data, nil
}

func getForm(name string) form {
	prefix := fmt.Sprintf("FORM_%s_", strings.ToUpper(name))
	form := form{
		name:    name,
		to:      os.Getenv(prefix + "TO"),
		from:    os.Getenv(prefix + "FROM"),
		subject: os.Getenv(prefix + "SUBJECT"),
	}
	return form
}

func formatData(form form, data map[string]string) string {
	message := fmt.Sprintf("<h1>New %s Submission</h1>", strings.Title(form.name))

	message += "<table><tbody>"
	for k, v := range data {
		if k[0] != '_' {
			message += fmt.Sprintf("<tr><th>%s</th><td>%s</td></tr>", k, v)
		}
	}
	message += "</tbody></table>"

	return message
}

func generateMessage(form form, message string) (string, error) {
	t, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return "", err
	}

	var email bytes.Buffer
	err = t.Execute(&email, message)
	if err != nil {
		return "", err
	}

	return inliner.Inline(email.String())
}

func sendEmail(server *mail.SMTPServer, form form, message string) error {
	email := mail.NewMSG()
	email.AddTo(form.to)
	email.SetFrom(form.from)
	email.SetSubject(form.subject)
	email.SetBody(mail.TextHTML, message)

	client, err := server.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	return email.Send(client)
}

// Handler is exported
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	server, err := setup()
	if err != nil {
		return respond(500, nil), nil
	}

	data, err := parseData(request.Headers["content-type"], request.Body)
	if err != nil {
		return respond(http.StatusBadRequest, err), nil
	}

	form := getForm(data["_form_name"])
	message := formatData(form, data)
	message, err = generateMessage(form, message)
	if err != nil {
		return respond(http.StatusInternalServerError, err), nil
	}

	err = sendEmail(server, form, message)
	if err != nil {
		return respond(http.StatusInternalServerError, err), nil
	}

	return respond(http.StatusOK, nil), nil
}
