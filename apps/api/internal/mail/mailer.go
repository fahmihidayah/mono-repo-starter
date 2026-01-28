package mail

import (
	"github.com/fahmihidayah/go-api-orchestrator/internal/config"
	"github.com/fahmihidayah/go-api-orchestrator/internal/domain"
	"github.com/go-gomail/gomail"
)

type Mailer struct {
	dialer         *gomail.Dialer
	templateEngine *TemplateEngine
	from           string
	frontendURL    string
}

func MailerProvider(cfg *config.Config, templateEngine *TemplateEngine) *Mailer {
	return &Mailer{
		dialer: gomail.NewDialer(
			cfg.Mail.Host,
			cfg.Mail.Port,
			cfg.Mail.User,
			cfg.Mail.Pass,
		),
		from:           cfg.Mail.From,
		templateEngine: templateEngine,
		frontendURL:    cfg.FrontendURL,
	}
}

func (m *Mailer) SendTo(to, subject, html string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", html)
	return m.dialer.DialAndSend(msg)
}

func (m *Mailer) SendResetPassword(user *domain.User) error {
	contentMessage, err := m.templateEngine.Render("reset_password.html", map[string]interface{}{
		"Name":     user.Name,
		"Email":    user.Email,
		"Token":    user.ResetPasswordToken,
		"ResetURL": m.frontendURL + "/reset-password?token=" + user.ResetPasswordToken,
	})
	if err != nil {
		return err
	}
	return m.SendTo(user.Email, "Reset Password", contentMessage)
}

func (m *Mailer) SendVerifyEmail(user *domain.User) error {
	contentMessage, err := m.templateEngine.Render("verify_email.html", map[string]interface{}{
		"Name":       user.Name,
		"Email":      user.Email,
		"VerifyLink": m.frontendURL + "/verify?email=" + user.Email + "&token=" + user.VerificationCode,
	})
	if err != nil {
		return err
	}
	return m.SendTo(user.Email, "Verify Link", contentMessage)
}
