# Getting Started with Google Cloud Functions

1. Create a Google Cloud Function 2nd generation.
2. Create a `function.go` in your project root:
```go
package gcf

import (
	"net/http"
	"github.com/djatwood/formailer"
	"github.com/djatwood/formailer/handlers"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
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

// Google Cloud Function entry point defined as "main":
func init() {
	functions.HTTP("main", Formailer)
}
```

3. Add your SMTP settings in the Cloud Functions UI.

4. Add a hidden input to your form.
```html
<input type="hidden" name="_form_name" value="Contact">
```
