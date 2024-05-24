
# Getting Started with Vercel

1. Create a `formailer.go` in your project at `./api/`:
```go
package api

import (
	"net/http"

	"github.com/torrayne/formailer"
	"github.com/torrayne/formailer/handlers"
)

// Formailer handles all form submissions
func Formailer(w http.ResponseWriter, r *http.Request) {
	contact := formailer.New("Contact")
	contact.AddEmail(formailer.Email{
		ID:      "contact",
		To:      "info@domain.com",
		From:    `"Company" <noreply@domain.com>`,
		Subject: "New Contact Submission",
	})

	handlers.Vercel(formailer.DefaultConfig, w, r)
}
```

2. Add your SMTP settings in you Vercel UI.
3. Add a hidden input to your form.
```html
<input type="hidden" name="_form_name" value="Contact">
```