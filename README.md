# Formailer

![Go](https://github.com/torrayne/formailer/workflows/Go/badge.svg)

![Screenshot](img.png)

If you need your contact form to send you an email from your Jamstack site, Formailer is the serverless library for you! Out of the box Formailer supports redirects, reCAPTCHA, custom email templates, and custom handlers.

## Quickstart
[View Documenation](https://pkg.go.dev/github.com/torrayne/formailer)

Formailer tries to require as little boilerplate as possible. Create a form, add some emails, and run a handler.
```go
import (
	"github.com/torrayne/formailer"
	"github.com/torrayne/formailer/handlers"
	
	// For Netlify
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	contact := formailer.New("Contact")
	contact.AddEmail(formailer.Email{
		To:      "info@domain.com",
		From:    `"Company" <noreply@domain.com>`,
		Subject: "New Contact Submission",
	})

	// Vercel
	handlers.Vercel(formailer.DefaultConfig, w, r)
	// Netlify
	lambda.Start(handlers.Netlify(formailer.DefaultConfig))
}
```
If you want to use your own handler that's not a problem either. [View an example handler](#user-content-custom-handlers).

Then you update your form. So that we can use the correct form config you need to add the form id as a hidden input with the name `_form_name`. Note the form id is case-insensitive.

Formailer supports submitting forms as `application/x-www-form-urlencoded`, `multipart/form-data`, or `application/json`.

The built-in handlers come with built in Google reCaptcha verification if you add a `RECAPTCHA_SECRET` to your environment variables.
```html
<!-- html form -->
<input type="hidden" name="_form_name" value="contact">
<button class="g-recaptcha" data-sitekey="reCAPTCHA_site_key" data-callback='onSubmit' data-action='submit'>Submit</button>
```
```javascript
// JSON object
{
	...
	"_form_name": "contact",
}
```

## Customization

You can customize Formailer to suit your needs. You can add as many forms as you'd like. As long as they have unique ids. Each form can have it's own email template and SMTP settings. But if you want to set defaults for everything you can.
### SMTP

All of your SMTP variables must be saved in the environment. You can add as many configs as you have emails. And you can save a default config to fallback on. Note that if you have default config you don't need to specify every option again. Any missing options will fallback to the default.

If you build your own handler you can store the config anywhere you want. Just pass a `*mail.SMTPServer` to `submission.Send(server)` and you're good to go.

```env
# Default
SMTP_HOST=mail.example.com
SMTP_PORT=587
SMTP_USER=noreply@example.com
SMTP_PASS=mysupersecretpassword

# Overrides
# _HOST and _PORT will fallback to the default above
SMTP_EMAIL-ID_USER=support@example.com
SMTP_EMAIL-ID_PASS=youcantguessthispassword
```

### Templates
Here is the default template.

![Screenshot](img.png)

You can override this template on any form by using the `Template` field. You can use Go 1.16 >= embed package to separate your template files from your function file.
```go
contact.AddEmail(formailer.Email{
	...
	Template: defaultTemplate,
}

//go:embed mytemplate.html
var defaultTemplate string

// OR

defaultTemplate := `
<html>
<head>
    <style>
        h3 {
            color: #000;
            margin-bottom: 5px;
        }

        p {
            color: #333;
            margin-bottom: 2rem;
        }
    </style>
</head>
<body>
    {{ range $name, $value := .Values }}
    <h3>{{$name}}</h3>
    <p>{{$value}}</p>
    {{ end }}
</body>
</html>
`
```

### Custom Handlers
Formailer ships with Netlify and Vercel handlers but if you need more control over the data. Or would like to run on a different platform, it's not too difficult to get setup. Here is a template to get you started.
```go
func Handler(w http.ResponseWriter, r *http.Request) {
	// pre-processing, check HTTP method

	// convert body from io.Reader to string
	body := new(strings.Builder)
	_, err := io.Copy(body, r.Body)
	if err != nil {
		// handle error
		return
	}
	
	// Parse body
	submission, err := formailer.Parse(r.Header.Get("Content-Type"), body.String())
	if err != nil {
		// handle error
		return
	}

	// manipulate data, check honey pot fields
	// handlers.VerifyRecaptcha()

	// Send emails
	err = submission.Send()
	if err != nil {
		// handle error
		return
	}

	// handle success
}
```

## Why did I buid Formailer?

I love Jamstack but SaaS can get expensive pretty quickly. Netlify has a built in form system that costs $19/month after the first 100 submissions. It also has a serverless function system that allows for 125k invocations a month. So I did the obvious thing, create a library that handles forms for Jamstack sites.

### The challenge
Netlify barely supports Go, you can't even use the Netlify CLI to test Go functions. Every change had to be commited and tested directly on Netlify. Even worse is that I had minimal experience working with multipart forms before this project. And my testing software [Hoppscotch](https://hoppscotch.io) doesn't implement multipart forms in a traditional way which led to a bunch of builds that I thought didn't work but actually did.

There's also an annoying bug with environment variables where [functions can't read variabes defined in the `netlify.toml`](https://github.com/netlify/netlify-lambda/issues/59). So you'll just have to add them all in the Netlify UI.

Later I switched to Vercel for testing which was a huge breath of fresh air. You can test Go functions locally even though Go support is still in alpha.
