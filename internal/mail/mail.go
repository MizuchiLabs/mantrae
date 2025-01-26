package mail

import (
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/MizuchiLabs/mantrae/internal/app"
	"github.com/MizuchiLabs/mantrae/internal/mail/templates"
	"github.com/domodwyer/mailyak/v3"
)

func Send(to, templateName string, email app.EmailConfig, data map[string]interface{}) error {
	client := mailyak.New(
		email.Host+":"+email.Port,
		smtp.PlainAuth("", email.Username, email.Password, email.Host),
	)
	client.To(to)
	client.From(email.From)

	var subject string
	switch templateName {
	case "reset-password":
		subject = "Reset your password"
	case "verify-email":
		subject = "Verify your email"
	default:
		return fmt.Errorf("unknown template: %s", templateName)
	}
	client.Subject(subject)

	// Parse the HTML template
	tmpl, err := template.New(templateName).
		ParseFS(templates.TemplateFS, templateName+".html")
	if err != nil {
		return err
	}

	// Output the result to the email client
	if err := tmpl.ExecuteTemplate(client.HTML(), templateName+".html", data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	if err := client.Send(); err != nil {
		return fmt.Errorf("failed to send email: %w, check your SMTP settings", err)
	}
	return nil
}
