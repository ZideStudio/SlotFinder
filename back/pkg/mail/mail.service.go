package mail

import (
	"app/commons/constants"
	"app/config"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"maps"
	"net/smtp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

//go:embed templates/*.html
var templateFS embed.FS

//go:embed locales/*/*.json
var localeFS embed.FS

type MailService struct {
	Config       config.Config
	templates    map[constants.MailTemplate]*template.Template
	translations map[string]map[constants.MailTemplate]map[string]any
}

func NewMailService(service *MailService) *MailService {
	if service != nil {
		return service
	}

	mailService := &MailService{
		Config:       *config.GetConfig(),
		templates:    make(map[constants.MailTemplate]*template.Template),
		translations: make(map[string]map[constants.MailTemplate]map[string]interface{}),
	}

	// Load all email templates
	if err := mailService.loadTemplates(); err != nil {
		log.Error().Err(err).Msg("MAIL_SERVICE::NEW Failed to load email templates")
	}

	// Load all translations
	if err := mailService.loadTranslations(); err != nil {
		log.Error().Err(err).Msg("MAIL_SERVICE::NEW Failed to load translations")
	}

	return mailService
}

type EmailParams struct {
	Template constants.MailTemplate
	To       string
	Subject  string
	Params   map[string]string
	Language constants.AccountLanguage
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

// loadTranslations loads all translation files from the locales directory
func (s *MailService) loadTranslations() error {
	// Read language directories
	languageDirs, err := localeFS.ReadDir("locales")
	if err != nil {
		log.Error().Err(err).Msg("MAIL_SERVICE::LOAD_TRANSLATIONS Failed to read locales directory")
		return fmt.Errorf("failed to read locales directory: %w", err)
	}

	for _, langDir := range languageDirs {
		if !langDir.IsDir() {
			continue
		}

		language := langDir.Name()
		languagePath := "locales/" + language

		// Read translation files in this language directory
		translationFiles, err := localeFS.ReadDir(languagePath)
		if err != nil {
			log.Error().Err(err).Str("language", language).Msg("MAIL_SERVICE::LOAD_TRANSLATIONS Failed to read language directory")
			continue
		}

		for _, file := range translationFiles {
			if !strings.HasSuffix(file.Name(), ".json") {
				continue
			}

			templateName := strings.TrimSuffix(file.Name(), ".json")
			translationPath := languagePath + "/" + file.Name()

			translationContent, err := localeFS.ReadFile(translationPath)
			if err != nil {
				log.Error().Err(err).Str("file", translationPath).Msg("MAIL_SERVICE::LOAD_TRANSLATIONS Failed to read translation file")
				continue
			}

			var translations map[string]any
			if err := json.Unmarshal(translationContent, &translations); err != nil {
				log.Error().Err(err).Str("file", translationPath).Msg("MAIL_SERVICE::LOAD_TRANSLATIONS Failed to parse translation file")
				continue
			}

			// Initialize nested maps if they don't exist
			if s.translations[language] == nil {
				s.translations[language] = make(map[constants.MailTemplate]map[string]any)
			}

			s.translations[language][constants.MailTemplate(templateName)] = translations
		}
	}

	return nil
}

// getTranslations returns translations for a specific template and language
func (s *MailService) getTranslations(templateName constants.MailTemplate, language constants.AccountLanguage) map[string]any {
	if langTranslations, exists := s.translations[string(language)]; exists {
		if templateTranslations, exists := langTranslations[templateName]; exists {
			return templateTranslations
		}
	}

	// Fallback to English if translation not found
	if langTranslations, exists := s.translations[string(constants.ACCOUNT_LANGUAGE_EN)]; exists {
		if templateTranslations, exists := langTranslations[templateName]; exists {
			log.Warn().
				Str("template", string(templateName)).
				Str("requested_language", string(language)).
				Msg("MAIL_SERVICE::GET_TRANSLATIONS Translation not found, falling back to English")
			return templateTranslations
		}
	}

	log.Error().
		Str("template", string(templateName)).
		Str("language", string(language)).
		Msg("MAIL_SERVICE::GET_TRANSLATIONS No translations found")
	return make(map[string]any)
}

// renderTemplate renders the email template with the provided parameters and translations
func (s *MailService) renderTemplate(templateName constants.MailTemplate, params map[string]string, language constants.AccountLanguage) (string, error) {
	tmpl, exists := s.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template '%s' not found", templateName)
	}

	// Get translations for the template and language
	translations := s.getTranslations(templateName, language)

	// Add translations
	templateData := make(map[string]any)
	maps.Copy(templateData, translations)

	// Add custom params (override case)
	for key, value := range params {
		templateData[key] = value
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
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

	// Set default language if not provided
	if params.Language == "" {
		params.Language = constants.ACCOUNT_LANGUAGE_EN
	}

	// Ensure params map is not nil
	if params.Params == nil {
		params.Params = make(map[string]string)
	}
	params.Params["Origin"] = s.Config.Origin

	// Render the email template with translations
	htmlBody, err := s.renderTemplate(params.Template, params.Params, params.Language)
	if err != nil {
		log.Error().
			Err(err).
			Str("template", string(params.Template)).
			Str("to", params.To).
			Str("language", string(params.Language)).
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
			Str("language", string(params.Language)).
			Msg("MAIL_SERVICE::SEND_MAIL Failed to send email")
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
