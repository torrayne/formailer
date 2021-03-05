
# Getting Started with Vercel

1. Create a `formailer.go` in your project at `./api/`:
```go
package main

import (
    "net/http"
	"github.com/djatwood/formailer"
    "github.com/djatwood/formailer/handlers"
)

func Send(w http.ResponseWriter, r *http.Request) {
    cfg := make(formailer.Config)
	cfg.Set(&formailer.Form{
        To:       "support@domain.com",
        From:     `"Company" <noreply@domain.com>`,
        Subject:  "New Submission",
        Redirect: "/success",
    })


	handlers.Vercel(&cfg, w, r)
}
```

2. Add your SMTP settings in you Vercel UI.
3. Add a hidden input to your form.
```html
<input type="hidden" name="_form_name" value="contact">
```