package mail

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"ksadebiekorf.be/mailing/db"
)

type Email = mail.Email

func NewEmail(name, address string) *Email {
	return mail.NewEmail(name, address)
}

// Parse email in the form of "name <email>"
func ParseEmail(e string) *Email {
	if len(e) == 0 {
		return nil
	}
	split := strings.Split(e, "<")
	name := strings.Trim(split[0], " ")
	email := strings.Trim(split[1], " \n>")
	// email := rhs[:len(rhs)-1]
	return NewEmail(name, email)
}

type MailConfig struct {
	From        *Email
	ReplyTo     *Email
	Subject     string
	Tos         []db.EmailEntry
	Content     string
	ContentType string
}

// Example: https://github.com/sendgrid/sendgrid-go/blob/bdbdc219ff7f51bd47087383c45002d43a877f07/examples/helpers/mail/example.go#L32
func CreateMailRequest(
	cfg *MailConfig,
) []byte {
	m := mail.NewV3Mail()
	p := mail.NewPersonalization()

	// From
	m.SetFrom(cfg.From)

	// Reply To
	m.SetReplyTo(cfg.ReplyTo)

	// Subject
	m.Subject = cfg.Subject

	// To
	var tos []*mail.Email
	for _, to := range cfg.Tos {
		tos = append(tos, mail.NewEmail(fmt.Sprintf("%s %s", to.FirstName, to.LastName), to.Email))
	}
	p.AddTos(tos...)

	// TODO: SetSubstitution, SetCustomArg
	// TODO: attachements

	content := mail.NewContent("text/html", cfg.HTMLContent)
	m.AddContent(content)

	// mailSettings := mail.NewMailSettings()
	// mailSettings.SetBypassSpamManagement(mail.NewSetting(true))
	m.AddPersonalizations(p)

	// TODO: footers

	return mail.GetRequestBody(m)
}

func SendEmail(
	requestBody []byte,
) error {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = requestBody
	resp, err := sendgrid.API(request)
	if err != nil {
		return err
	}

	log.Printf("Email sent: [%d] %s\n", resp.StatusCode, resp.Body)
	return nil
}
