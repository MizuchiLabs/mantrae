package util

import (
	"embed"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/domodwyer/mailyak/v3"
)

//go:embed templates/*.html
var templatesFS embed.FS

func SendMail(to, templateName string, config EmailConfig, data map[string]interface{}) error {
	client := mailyak.New(
		config.EmailHost+":"+config.EmailPort,
		smtp.PlainAuth("", config.EmailUsername, config.EmailPassword, config.EmailHost),
	)
	client.To(to)
	client.From(config.EmailFrom)

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
	tmpl, err := template.New(templateName).ParseFS(templatesFS, "templates/"+templateName+".html")
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
