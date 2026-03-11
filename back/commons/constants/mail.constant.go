package constants

type MailTemplate string

const (
	MAIL_TEMPLATE_WELCOME                     MailTemplate = "welcome"
	MAIL_TEMPLATE_PASSWORD_RESET              MailTemplate = "password-reset"
	MAIL_TEMPLATE_PASSWORD_RESET_CONFIRMATION MailTemplate = "password-reset-confirmation"
)

const (
	MAIL_SUBJECT_WELCOME_EN                = "Welcome to SlotFinder!"
	MAIL_SUBJECT_WELCOME_FR                = "Bienvenue sur SlotFinder !"
	MAIL_SUBJECT_PASSWORD_RESET_EN         = "Reset your password"
	MAIL_SUBJECT_PASSWORD_RESET_FR         = "Réinitialiser votre mot de passe"
	MAIL_SUBJECT_PASSWORD_RESET_CONFIRM_EN = "Password reset successful"
	MAIL_SUBJECT_PASSWORD_RESET_CONFIRM_FR = "Mot de passe réinitialisé avec succès"
)
