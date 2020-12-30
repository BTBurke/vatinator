package vatinator

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/keighl/postmark"
)

const textEmail = `Hi {{.FormData.FirstName}},

Your forms for {{.Month}} are ready.  Click or paste the link below to download them:

{{.Link}}

If you notice any problems, you can reply to this email for help.
`

type EmailData struct {
	FormData FormData
	Month    string
	Year     int
	Link     string
}

type EmailService interface {
	SendDownloadEmail(address string, data EmailData) error
}

type emailService struct {
	client *postmark.Client
}

func NewEmailService(serverToken, apiToken string) EmailService {
	return emailService{postmark.NewClient(serverToken, apiToken)}
}

func (e emailService) SendDownloadEmail(address string, data EmailData) error {
	t, err := template.New("text").Parse(textEmail)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	if err := t.Execute(&b, data); err != nil {
		return err
	}

	email := postmark.Email{
		From:     "forms@vatinator.com",
		To:       address,
		Subject:  fmt.Sprintf("Your VAT forms for %s %d", data.Month, data.Year),
		HtmlBody: "",
		TextBody: b.String(),
	}

	_, err = e.client.SendEmail(email)
	if err != nil {
		return err
	}
	return nil
}
