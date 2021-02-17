package formailer

//go:generate go run generate/main.go

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
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
	name     string
	to       string
	from     string
	subject  string
	redirect string
}

type formData struct {
	values      map[string]string
	attachments []attachment
}

type attachment struct {
	filename string
	mimeType string
	data     []byte
}

func setup() (*mail.SMTPServer, error) {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return nil, errors.New("could not parse SMTP_PORT")
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

func vercelResponse(w http.ResponseWriter, code int, err error, headers ...[2]string) {
	for _, h := range headers {
		w.Header().Add(h[0], h[1])
	}
	body := http.StatusText(code)
	if err != nil {
		body = err.Error()
		w.Header().Set("location", w.Header().Get("location")+"?error="+err.Error())
	}

	w.WriteHeader(code)
	w.Write([]byte(body))
}

func readerToString(src io.Reader) (string, error) {
	dst := new(strings.Builder)
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return dst.String(), nil
}

func parseData(contentType string, body string) (*formData, error) {
	data := new(formData)
	data.values = make(map[string]string)

	var err error
	if strings.Contains(contentType, "application/json") {
		err = data.parseJSON(body)
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err = data.parseURLEncoded(body)
	} else if strings.Contains(contentType, "multipart/form-data") {
		err = data.parseMultipartForm(contentType, body)
	} else {
		err = errors.New("invalid content type")
	}

	return data, err
}

func (data *formData) parseJSON(body string) error {
	return json.Unmarshal([]byte(body), &data.values)
}

func (data *formData) parseURLEncoded(body string) error {
	vals, err := url.ParseQuery(body)
	if err != nil {
		return err
	}

	for k := range vals {
		data.values[k] = vals.Get(k)
	}
	return nil
}

func (data *formData) parseMultipartForm(contentType, body string) error {
	header := strings.Split(contentType, ";")
	var boundary string
	for _, h := range header {
		index := strings.Index(h, "boundary")
		if index > -1 {
			boundary = strings.TrimSpace(h[index+9:])
			break
		}
	}

	decodedBody, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		return err
	}
	decodedBody = append(decodedBody, '\n')

	reader := multipart.NewReader(bytes.NewReader(decodedBody), boundary)
	for {
		part, err := reader.NextPart()
		if err != nil {
			return err
		}

		value := new(bytes.Buffer)
		value.ReadFrom(part)

		filename := part.FileName()
		if len(filename) > 0 {
			attachment := attachment{
				filename: part.FileName(),
				mimeType: part.Header.Get("Content-Type"),
				data:     value.Bytes(),
			}
			data.attachments = append(data.attachments, attachment)
			data.values[part.FormName()] = part.FileName()
		} else {
			data.values[part.FormName()] = value.String()
		}
	}
}

func (data *formData) format(form form) string {
	message := fmt.Sprintf("<h1>New %s Submission</h1>", strings.Title(form.name))

	message += "<table><tbody>"
	for k, v := range data.values {
		if k[0] != '_' {
			message += fmt.Sprintf("<tr><th>%s</th><td>%s</td></tr>", k, v)
		}
	}
	message += "</tbody></table>"

	return message
}

func getForm(name string) (form, error) {
	if name == "" {
		return form{}, errors.New("_form_name missing from input")
	}

	prefix := fmt.Sprintf("FORM_%s_", strings.ToUpper(name))
	form := form{
		name:     name,
		to:       os.Getenv(prefix + "TO"),
		from:     os.Getenv(prefix + "FROM"),
		subject:  os.Getenv(prefix + "SUBJECT"),
		redirect: os.Getenv(prefix + "REDIRECT"),
	}

	err := errors.New("could not parse form from env: missing")
	if form.to == "" {
		return form, fmt.Errorf("%w to", err)
	}
	if form.from == "" {
		return form, fmt.Errorf("%w from", err)
	}
	if form.subject == "" {
		return form, fmt.Errorf("%w subject", err)
	}

	return form, nil
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

func sendEmail(server *mail.SMTPServer, form form, data *formData) error {
	message := data.format(form)
	message, err := generateMessage(form, message)
	if err != nil {
		return err
	}

	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return errors.New("failed to generate message-id")
	}

	email := mail.NewMSG()
	email.AddTo(form.to)
	email.SetFrom(form.from)
	email.SetSubject(form.subject)
	email.SetBody(mail.TextHTML, message)
	email.AddHeader("Message-Id", base32.StdEncoding.EncodeToString(token))

	for _, attachment := range data.attachments {
		email.AddAttachmentData(attachment.data, attachment.filename, attachment.mimeType)
	}

	client, err := server.Connect()
	if err != nil {
		return err
	}
	defer client.Close()
	return email.Send(client)
}

// Netlify takes in a aws lambda request and sends an email
func Netlify(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != "POST" {
		return netlifyResponse(http.StatusMethodNotAllowed, errors.New("only supports POST requests")), nil
	}

	server, err := setup()
	if err != nil {
		return netlifyResponse(500, nil), nil
	}

	data, err := parseData(request.Headers["content-type"], request.Body)
	if err != nil && err != io.EOF {
		return netlifyResponse(http.StatusBadRequest, err), nil
	}
	if v := data.values["faxonly"]; v == "1" {
		return netlifyResponse(http.StatusOK, nil), nil
	}

	form, err := getForm(data.values["_form_name"])
	if err != nil {
		return netlifyResponse(http.StatusBadRequest, err), nil
	}

	err = sendEmail(server, form, data)
	if err != nil {
		return netlifyResponse(http.StatusInternalServerError, err), nil
	}

	statusCode := http.StatusOK
	headers := [][2]string{}
	if len(form.redirect) > 0 {
		statusCode = http.StatusSeeOther
		headers = append(headers, [2]string{"location", form.redirect})
	}

	return netlifyResponse(statusCode, nil, headers...), nil
}

// Vercel just needs a normal http handler
func Vercel(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		vercelResponse(w, http.StatusMethodNotAllowed, errors.New("only supports POST requests"))
		return
	}

	server, err := setup()
	if err != nil {
		vercelResponse(w, http.StatusInternalServerError, err)
		return
	}

	body := new(strings.Builder)
	_, err = io.Copy(body, r.Body)
	if err != nil {
		vercelResponse(w, http.StatusInternalServerError, err)
		return
	}
	data, err := parseData(r.Header.Get("Content-Type"), body.String())
	if err != nil && err != io.EOF {
		vercelResponse(w, http.StatusBadRequest, err)
		return
	}
	if v := data.values["faxonly"]; v == "1" {
		vercelResponse(w, http.StatusOK, nil)
		return
	}

	form, err := getForm(data.values["_form_name"])
	if err != nil {
		vercelResponse(w, http.StatusBadRequest, err)
		return
	}

	err = sendEmail(server, form, data)
	if err != nil {
		vercelResponse(w, http.StatusInternalServerError, err)
		return
	}

	statusCode := http.StatusOK
	headers := [][2]string{}
	if len(form.redirect) > 0 {
		statusCode = http.StatusSeeOther
		headers = append(headers, [2]string{"location", form.redirect})
	}

	vercelResponse(w, statusCode, nil, headers...)
}
