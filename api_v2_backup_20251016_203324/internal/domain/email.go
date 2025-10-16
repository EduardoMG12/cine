package domain

// EmailService defines the interface for email operations
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
	SendConfirmationEmail(to, username, confirmationURL string) error
	SendPasswordResetEmail(to, username, resetURL string) error
	SendSecureEmail(to, subject, body string) error
	TestConnection() error
}

// EmailTemplate represents different email template types
type EmailTemplate string

const (
	EmailTemplateConfirmation  EmailTemplate = "confirmation"
	EmailTemplatePasswordReset EmailTemplate = "password_reset"
	EmailTemplateWelcome       EmailTemplate = "welcome"
	EmailTemplateNotification  EmailTemplate = "notification"
)

// EmailData represents the data structure for email templates
type EmailTemplateData struct {
	Username        string
	Email           string
	ConfirmationURL string
	ResetURL        string
	AppName         string
	SupportEmail    string
	ExpiresIn       string
	// Additional fields can be added as needed
	CustomData map[string]interface{}
}
