package constants

type MailTemplate string

const (
	MAIL_TEMPLATE_WELCOME                     MailTemplate = "welcome"
	MAIL_TEMPLATE_PASSWORD_RESET              MailTemplate = "password-reset"
	MAIL_TEMPLATE_PASSWORD_RESET_CONFIRMATION MailTemplate = "password-reset-confirmation"
)
