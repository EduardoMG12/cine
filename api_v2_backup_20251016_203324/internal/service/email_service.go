package service

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/EduardoMG12/cine/api_v2/internal/config"
)

type EmailService struct {
	cfg *config.EmailConfig
}

func NewEmailService(cfg *config.EmailConfig) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

// EmailData represents data to be passed to email templates
type EmailData struct {
	Username        string
	ConfirmationURL string
	ResetURL        string
	AppName         string
	SupportEmail    string
}

// SendEmail sends a basic email
func (s *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	msg := s.buildMessage(to, subject, body, false)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, msg)
}

// SendHTMLEmail sends an HTML email
func (s *EmailService) SendHTMLEmail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)

	msg := s.buildMessage(to, subject, htmlBody, true)

	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)
	return smtp.SendMail(addr, auth, s.cfg.FromEmail, []string{to}, msg)
}

// SendConfirmationEmail sends an email confirmation email
func (s *EmailService) SendConfirmationEmail(to, username, confirmationURL string) error {
	data := EmailData{
		Username:        username,
		ConfirmationURL: confirmationURL,
		AppName:         "CineVerse",
		SupportEmail:    s.cfg.FromEmail,
	}

	htmlBody, err := s.renderTemplate("confirmation", data)
	if err != nil {
		return fmt.Errorf("failed to render confirmation template: %w", err)
	}

	subject := "Welcome to CineVerse - Please Confirm Your Email"
	return s.SendHTMLEmail(to, subject, htmlBody)
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, username, resetURL string) error {
	data := EmailData{
		Username:     username,
		ResetURL:     resetURL,
		AppName:      "CineVerse",
		SupportEmail: s.cfg.FromEmail,
	}

	htmlBody, err := s.renderTemplate("password_reset", data)
	if err != nil {
		return fmt.Errorf("failed to render password reset template: %w", err)
	}

	subject := "CineVerse - Password Reset Request"
	return s.SendHTMLEmail(to, subject, htmlBody)
}

// SendSecureEmail sends email with TLS encryption
func (s *EmailService) SendSecureEmail(to, subject, body string) error {
	// Create TLS config
	tlsConfig := &tls.Config{
		ServerName: s.cfg.SMTPHost,
	}

	// Connect to the SMTP server with TLS
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with TLS: %w", err)
	}
	defer conn.Close()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.cfg.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// Authenticate
	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender and recipient
	if err := client.Mail(s.cfg.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send the email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	msg := s.buildMessage(to, subject, body, true)
	if _, err := writer.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return writer.Close()
}

// buildMessage constructs the email message
func (s *EmailService) buildMessage(to, subject, body string, isHTML bool) []byte {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", s.cfg.FromName, s.cfg.FromEmail))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	if isHTML {
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	msg.WriteString("\r\n")
	msg.WriteString(body)

	return []byte(msg.String())
}

// renderTemplate renders email templates
func (s *EmailService) renderTemplate(templateName string, data EmailData) (string, error) {
	var templateStr string

	switch templateName {
	case "confirmation":
		templateStr = confirmationEmailTemplate
	case "password_reset":
		templateStr = passwordResetEmailTemplate
	default:
		return "", fmt.Errorf("unknown template: %s", templateName)
	}

	tmpl, err := template.New(templateName).Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// TestConnection tests the SMTP connection
func (s *EmailService) TestConnection() error {
	slog.Info("Testing SMTP connection", "host", s.cfg.SMTPHost, "port", s.cfg.SMTPPort)

	auth := smtp.PlainAuth("", s.cfg.SMTPUsername, s.cfg.SMTPPassword, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	// Try to connect and authenticate
	conn, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Quit()

	if err := conn.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP server: %w", err)
	}

	slog.Info("SMTP connection test successful")
	return nil
}

// Email templates
const confirmationEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to {{.AppName}}</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #1f2937; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #3b82f6; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; font-size: 14px; color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to {{.AppName}}!</h1>
        </div>
        <div class="content">
            <h2>Hi {{.Username}},</h2>
            <p>Welcome to {{.AppName}}! We're excited to have you join our community of movie enthusiasts.</p>
            <p>To complete your registration and start discovering amazing movies, please confirm your email address by clicking the button below:</p>
            <p style="text-align: center;">
                <a href="{{.ConfirmationURL}}" class="button">Confirm Your Email</a>
            </p>
            <p>If the button above doesn't work, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #3b82f6;">{{.ConfirmationURL}}</p>
            <p>Once confirmed, you'll be able to:</p>
            <ul>
                <li>Rate and review movies</li>
                <li>Create custom movie lists</li>
                <li>Follow other users and see their recommendations</li>
                <li>Discover personalized movie recommendations</li>
            </ul>
        </div>
        <div class="footer">
            <p>If you didn't create an account with {{.AppName}}, please ignore this email.</p>
            <p>Need help? Contact us at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a></p>
        </div>
    </div>
</body>
</html>
`

const passwordResetEmailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your {{.AppName}} Password</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { background-color: #dc2626; color: white; padding: 20px; text-align: center; border-radius: 8px 8px 0 0; }
        .content { background-color: #f9fafb; padding: 30px; border-radius: 0 0 8px 8px; }
        .button { display: inline-block; padding: 12px 24px; background-color: #dc2626; color: white; text-decoration: none; border-radius: 6px; margin: 20px 0; }
        .footer { text-align: center; margin-top: 20px; font-size: 14px; color: #666; }
        .warning { background-color: #fef3c7; border-left: 4px solid #f59e0b; padding: 12px; margin: 16px 0; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Password Reset Request</h1>
        </div>
        <div class="content">
            <h2>Hi {{.Username}},</h2>
            <p>We received a request to reset your {{.AppName}} password. If you made this request, click the button below to reset your password:</p>
            <p style="text-align: center;">
                <a href="{{.ResetURL}}" class="button">Reset Your Password</a>
            </p>
            <p>If the button above doesn't work, you can copy and paste this link into your browser:</p>
            <p style="word-break: break-all; color: #dc2626;">{{.ResetURL}}</p>
            <div class="warning">
                <strong>Important:</strong> This password reset link will expire in 1 hour for security reasons.
            </div>
            <p>If you didn't request a password reset, please ignore this email. Your password will remain unchanged.</p>
        </div>
        <div class="footer">
            <p>For security reasons, this link will expire soon. If you need help, contact us at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a></p>
        </div>
    </div>
</body>
</html>
`
