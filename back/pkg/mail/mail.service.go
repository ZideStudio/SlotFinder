package mail

import (
	"app/commons/constants"
	"app/config"
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*.html
var templateFS embed.FS

type MailService struct {
	Config    config.Config
	templates map[constants.MailTemplate]*template.Template
}

func NewMailService(service *MailService) *MailService {
	if service != nil {
		return service
	}

	mailService := &MailService{
		Config:    *config.GetConfig(),
		templates: make(map[constants.MailTemplate]*template.Template),
	}

	// Load all email templates
	if err := mailService.loadTemplates(); err != nil {
		log.Error().Err(err).Msg("MAIL_SERVICE::NEW Failed to load email templates")
	}

	return mailService
}

type EmailParams struct {
	Template constants.MailTemplate
	To       string
	Subject  string
	Params   map[string]string
}

// loadTemplates loads all HTML templates from the templates directory
func (s *MailService) loadTemplates() error {
	templateFiles, err := templateFS.ReadDir("templates")
	if err != nil {
		log.Error().Err(err).Msg("MAIL_SERVICE::LOAD_TEMPLATES Failed to read templates directory")
		return fmt.Errorf("failed to read templates directory: %w", err)
	}

	for _, file := range templateFiles {
		if !strings.HasSuffix(file.Name(), ".html") {
			continue
		}

		templateName := strings.TrimSuffix(file.Name(), ".html")
		templatePath := "templates/" + file.Name()

		templateContent, err := templateFS.ReadFile(templatePath)
		if err != nil {
			log.Error().Err(err).Str("template", templateName).Msg("MAIL_SERVICE::LOAD_TEMPLATES Failed to read template file")
			continue
		}

		tmpl, err := template.New(templateName).Parse(string(templateContent))
		if err != nil {
			log.Error().Err(err).Str("template", templateName).Msg("MAIL_SERVICE::LOAD_TEMPLATES Failed to parse template")
			continue
		}

		s.templates[constants.MailTemplate(templateName)] = tmpl
	}

	return nil
}

// renderTemplate renders the email template with the provided parameters
func (s *MailService) renderTemplate(templateName constants.MailTemplate, params map[string]string) (string, error) {
	tmpl, exists := s.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template '%s' not found", templateName)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %w", templateName, err)
	}

	return buf.String(), nil
}

// buildEmailMessage constructs the complete email message
func (s *MailService) buildEmailMessage(to, subject, htmlBody string) []byte {
	messageID := fmt.Sprintf(
		"<%s@%s>",
		uuid.NewString(),
		strings.Split(s.Config.Email.Address, "@")[1],
	)

	date := time.Now().UTC().Format(time.RFC1123Z)

	message := fmt.Sprintf(
		"From: SlotFinder <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Date: %s\r\n"+
			"Message-ID: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"MIME-Version: 1.0\r\n"+
			"\r\n"+
			"%s\r\n",
		s.Config.Email.Address,
		to,
		subject,
		date,
		messageID,
		htmlBody,
	)

	return []byte(message)
}

// SendMail sends an email using the specified template and parameters
func (s *MailService) SendMail(params EmailParams) error {
	// Validate required parameters
	if params.Template == "" {
		return fmt.Errorf("template name is required")
	}
	if params.To == "" {
		return fmt.Errorf("recipient email is required")
	}
	if params.Subject == "" {
		return fmt.Errorf("email subject is required")
	}

	// Ensure params map is not nil
	if params.Params == nil {
		params.Params = make(map[string]string)
	}
	params.Params["Origin"] = s.Config.Origin

	// Render the email template
	htmlBody, err := s.renderTemplate(params.Template, params.Params)
	if err != nil {
		log.Error().
			Err(err).
			Str("template", string(params.Template)).
			Str("to", params.To).
			Msg("MAIL_SERVICE::SEND_MAIL Failed to render email template")
		return fmt.Errorf("failed to render email template: %w", err)
	}

	// Build the complete email message
	message := s.buildEmailMessage(params.To, params.Subject, htmlBody)

	// Configure SMTP authentication
	smtpHost := s.Config.Email.Host
	smtpPort := s.Config.Email.Port
	from := s.Config.Email.Address
	password := s.Config.Email.Password

	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{params.To},
		message,
	)

	if err != nil {
		log.Error().
			Err(err).
			Str("template", string(params.Template)).
			Str("to", params.To).
			Str("subject", params.Subject).
			Msg("MAIL_SERVICE::SEND_MAIL Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
