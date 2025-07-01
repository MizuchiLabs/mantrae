// Package mail provides functionality for sending emails.
package mail

import (
	"context"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/domodwyer/mailyak/v3"
	"github.com/mizuchilabs/mantrae/internal/mail/templates"
	"github.com/mizuchilabs/mantrae/internal/settings"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func Send(sm *settings.SettingsManager, to, templateName string, data map[string]any) error {
	email, err := getConfig(sm)
	if err != nil {
		return err
	}

	client := mailyak.New(
		email.Host+":"+email.Port,
		smtp.PlainAuth("", email.Username, email.Password, email.Host),
	)
	client.To(to)
	client.From(email.From)
	client.FromName("Mantrae")

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

func getConfig(sm *settings.SettingsManager) (*EmailConfig, error) {
	sets := sm.GetMany(context.Background(), []string{
		settings.KeyEmailHost,
		settings.KeyEmailPort,
		settings.KeyEmailUser,
		settings.KeyEmailPassword,
		settings.KeyEmailFrom,
	})

	return &EmailConfig{
		Host:     sets[settings.KeyEmailHost],
		Port:     sets[settings.KeyEmailPort],
		Username: sets[settings.KeyEmailUser],
		Password: sets[settings.KeyEmailPassword],
		From:     sets[settings.KeyEmailFrom],
	}, nil
}
