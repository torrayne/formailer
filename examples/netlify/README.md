# Getting Started with Netlify

1. Create a `main.go` in your project root:
```go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djatwood/formailer"
)

func main() {
    cfg := make(formialer.Config)
    cfg.Add(&formailer.Form{
        To:       "support@domain.com",
        From:     `"Company" <noreply@domain.com>`,
	    Subject:  "New Submission",
	    Redirect: "/success",
    })
	lambda.Start(cfg.Netlify)
}
```
2. Update your `netlify.toml`:
```toml
[build]
    build="go build -o functions/formailer"
    functions="functions" 
[build.environment]
    GO_IMPORT_PATH="your project git location"
```

3. Add your SMTP settings in you Netlify UI.

4. Add a hidden input to your form.
```html
<input type="hidden" name="_form_name" value="contact">
```