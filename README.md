# Formailer
![Go](https://github.com/djatwood/formailer/workflows/Go/badge.svg)
![Screenshot](img.png)

Netlify only allows for 100 form submissions a month but also allows for 125k function calls per month. So for websites that don't need submission storage you can use Formailer to send you an email after submission.

## Challenges

Netlify barely supports Go, you can't even use the Netlify CLI to test Go functions. Every change had to be commited and tested directly on Netlify. The bad part is I had minimal experience working with multipart forms before this project. And [Hoppscotch](https://hoppscotch.io) doesn't implement multipart forms in a traditional way which led to a bunch of builds that I thought didn't work but actually did.

There's also an annoying bug with environment variables where [functions can't read variabes defined in the `netlify.toml`](https://github.com/netlify/netlify-lambda/issues/59). So you'll just have to add them all in the UI.

## Install
Create a `main.go` in your project root:
```go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djatwood/formailer"
)

func main() {
	lambda.Start(formailer.Handler)
}
```
Update your `netlify.toml`:
```toml
[build]
    build="go build -o functions/formailer"
    functions="functions" 
[build.environment]
    GO_IMPORT_PATH="your project git location"
```
Here is [another example](https://github.com/netlify/aws-lambda-go-example)

## Setup
I wanted to keep your email secure. There are several form email services that leak your email either inside hidden fields or in the form action. And while that may be okay, I really didn't want to do it that way. So you add your SMTP config to your env.

### SMTP
```env
SMTP_HOST=mail.domain.com
SMTP_PORT=587
SMTP_USER=noreply@domain.com
SMTP_PASS=mysupersecretpassword
```

### Forms
Each form will need the following vars: `_TO`, `_FROM`, and `SUBJECT`. `_REDIRECT` is the location you want the user to go to after submitting the form. This is unnecessary for forms submitted with AJAX.
```env
FORM_FORMNAME_TO=
FORM_FORMNAME_FROM="Name" <email@domail.com>
FORM_FORMNAME_SUBJECT=
# OPTIONAL
FORM_FORMNAME_REDIRECT=
```

Your setup might look like this:
```env
FORM_QUOTES_TO=support@domain.com
FORM_QUOTES_FROM="Company" <noreply@domain.com>
FORM_QUOTES_SUBJECT=New Quote Request
FORM_REDIRECT=/quotes/success

FORM_CONTACT_TO=info@domain.com
FORM_CONTACT_FROM="Company" <noreply@domain.com>
FORM_CONTACT_SUBJECT=New Contact Submission
FORM_REDIRECT=https://domin.com/thankyou
```

### Submitting forms
So that we can use the correct form config you need to add the form name as a hidden input with the name `_form_name`. You can also add a honey pot checkbox with the name `faxonly`.

Formailer supports submitting forms as `application/x-www-form-urlencoded`, `multipart/form-data`, or `application/json`.
```html
<input name="_form_name" value="contact">
<!-- Honey Pot -->
<input type="checkbox" name="faxonly" value="1" style="display:none !important" tabindex="-1" autocomplete="off">
```

