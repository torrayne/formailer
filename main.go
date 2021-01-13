package main

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

const stylesheet = `
body {
	font-family: sans-serif;
}

#wrapper {
	padding: 30px 20px;
	border: 1px solid #ccc;
	border-radius: 4px;
	margin: 20px 0 10px;
}

.content {
	width: 100%;
	max-width: 800px;
	margin: 0;
	box-sizing: border-box;
}

h1 {
	font-size: 1.8rem;
	color: #212121;
	font-weight: bold;
	margin: 0 0 10px;
}

table {
	border-collapse: collapse;
}

th,
td {
	vertical-align: top;
	padding: 10px 10px;
	border-top: 1px solid #ddd;
	text-align: left;
	color: #404040;
	font-size: 1.1rem;
	font-weight: 400;
}

th {
	font-weight: bold;
}

.attribute {
	padding: 0 20px;
	text-decoration: none;
	color: #404040;
}
`

const template = `
<html lang="en">
<head><style>%s</style></head>
<body>
<div id="wrapper" class="content">%s</div>
<p class="content"><a class="attribute" href="https://atwood.io">Powered by Formailer © Atwood.io</a></p></body>
</html>`

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

func getData(request events.APIGatewayProxyRequest) (map[string]string, error) {
	data := make(map[string]string)
	contentType := request.Headers["content-type"]

	if strings.Contains(contentType, "application/json") {
		err := json.Unmarshal([]byte(request.Body), &data)
		if err != nil {
			return data, err
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		vals, err := url.ParseQuery(request.Body)
		if err != nil {
			return data, err
		}
		for k := range vals {
			data[k] = vals.Get(k)
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		fmt.Println(request.Body)
		reader := multipart.NewReader(strings.NewReader(request.Body), "\n")
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			buf := new(bytes.Buffer)
			buf.ReadFrom(part)
			fmt.Println(buf.String())
		}
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
	return inliner.Inline(fmt.Sprintf(template, stylesheet, message))
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

	data, err := getData(request)
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
