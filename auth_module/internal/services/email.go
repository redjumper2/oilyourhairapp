package services

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/sparque/auth_module/config"
	"github.com/sparque/auth_module/internal/models"
)

// EmailService handles sending emails
type EmailService struct {
	cfg *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{cfg: cfg}
}

// SendMagicLink sends a magic link email
func (s *EmailService) SendMagicLink(to, domain, magicLinkURL string) error {
	subject := fmt.Sprintf("Sign in to %s", domain)

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; }
		.footer { margin-top: 30px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		<h2>Sign in to %s</h2>
		<p>Click the button below to sign in to your account:</p>
		<p><a href="%s" class="button">Sign In</a></p>
		<p>Or copy and paste this link into your browser:</p>
		<p><a href="%s">%s</a></p>
		<p class="footer">
			This link will expire in 15 minutes.<br>
			If you didn't request this email, you can safely ignore it.
		</p>
	</div>
</body>
</html>
	`, domain, magicLinkURL, magicLinkURL, magicLinkURL)

	return s.sendEmail(to, subject, body)
}

// SendInvitation sends an invitation email with QR code
func (s *EmailService) SendInvitation(to, domain string, invitation *models.Invitation, inviteURL, qrCodeDataURL string, branding *models.DomainBranding) error {
	subject := fmt.Sprintf("You've been invited to %s", branding.CompanyName)

	tmpl := template.Must(template.New("invitation").Parse(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.button { display: inline-block; padding: 12px 24px; background-color: {{.PrimaryColor}}; color: white; text-decoration: none; border-radius: 4px; }
		.qr-code { margin: 20px 0; text-align: center; }
		.qr-code img { max-width: 256px; border: 1px solid #ddd; padding: 10px; }
		.footer { margin-top: 30px; font-size: 12px; color: #666; }
	</style>
</head>
<body>
	<div class="container">
		{{if .LogoURL}}
		<img src="{{.LogoURL}}" alt="{{.CompanyName}}" style="max-width: 200px;">
		{{end}}

		<h2>You've been invited to {{.CompanyName}}</h2>
		<p>You've been invited to join <strong>{{.CompanyName}}</strong> as a <strong>{{.Role}}</strong>.</p>

		<p><a href="{{.InviteURL}}" class="button">Accept Invitation</a></p>

		<p>Or scan this QR code:</p>
		<div class="qr-code">
			<img src="{{.QRCodeDataURL}}" alt="QR Code">
		</div>

		<p>Or copy and paste this link:</p>
		<p><a href="{{.InviteURL}}">{{.InviteURL}}</a></p>

		<p class="footer">
			This invitation expires on {{.ExpiresAt}}.<br>
			{{if .SupportEmail}}
			Questions? Contact us at <a href="mailto:{{.SupportEmail}}">{{.SupportEmail}}</a>
			{{end}}
		</p>
	</div>
</body>
</html>
	`))

	var body bytes.Buffer
	data := map[string]interface{}{
		"CompanyName":    branding.CompanyName,
		"PrimaryColor":   branding.PrimaryColor,
		"LogoURL":        branding.LogoURL,
		"Role":           invitation.Role,
		"InviteURL":      inviteURL,
		"QRCodeDataURL":  qrCodeDataURL,
		"ExpiresAt":      invitation.ExpiresAt.Format("January 2, 2006 at 3:04 PM"),
		"SupportEmail":   branding.SupportEmail,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to render email template: %w", err)
	}

	return s.sendEmail(to, subject, body.String())
}

// sendEmail sends an email using SMTP
func (s *EmailService) sendEmail(to, subject, body string) error {
	from := s.cfg.Email.FromAddress

	// Setup email headers
	msg := []byte(fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", from, to, subject, body))

	// SMTP authentication
	auth := smtp.PlainAuth("",
		s.cfg.Email.SMTP.User,
		s.cfg.Email.SMTP.Password,
		s.cfg.Email.SMTP.Host,
	)

	// Send email
	addr := fmt.Sprintf("%s:%d", s.cfg.Email.SMTP.Host, s.cfg.Email.SMTP.Port)
	err := smtp.SendMail(addr, auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
