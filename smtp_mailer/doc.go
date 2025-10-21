// Package smtp_mailer provides a simple SMTP email client with TLS support
// for sending HTML and plain text emails.
//
// # Basic Usage
//
// Create a mailer and send an email:
//
//	import "github.com/poly-workshop/go-webmods/smtp_mailer"
//
//	mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
//	    Host:      "smtp.gmail.com",
//	    Port:      587,
//	    Username:  "your-email@gmail.com",
//	    Password:  "your-app-password",
//	    FromEmail: "your-email@gmail.com",
//	    FromName:  "My Application",
//	})
//
//	err := mailer.SendEmail(smtp_mailer.Message{
//	    ToEmails: []string{"user@example.com"},
//	    Subject:  "Welcome to our service",
//	    Body:     "Thank you for signing up!",
//	    IsHTML:   false,
//	})
//
// # HTML Emails
//
// Send HTML-formatted emails:
//
//	err := mailer.SendEmail(smtp_mailer.Message{
//	    ToEmails: []string{"user@example.com"},
//	    Subject:  "Welcome!",
//	    Body: `
//	        <html>
//	            <body>
//	                <h1>Welcome to our service!</h1>
//	                <p>Thank you for signing up.</p>
//	                <a href="https://example.com/verify">Verify Email</a>
//	            </body>
//	        </html>
//	    `,
//	    IsHTML: true,
//	})
//
// # Multiple Recipients
//
// Send to multiple recipients:
//
//	err := mailer.SendEmail(smtp_mailer.Message{
//	    ToEmails: []string{
//	        "user1@example.com",
//	        "user2@example.com",
//	        "user3@example.com",
//	    },
//	    Subject: "Newsletter",
//	    Body:    "Check out our latest updates!",
//	})
//
// # Common SMTP Providers
//
// Gmail (requires App Password):
//
//	Config{
//	    Host:     "smtp.gmail.com",
//	    Port:     587,
//	    Username: "your-email@gmail.com",
//	    Password: "app-password",  // Not your regular password
//	}
//
// Office 365:
//
//	Config{
//	    Host:     "smtp.office365.com",
//	    Port:     587,
//	    Username: "your-email@outlook.com",
//	    Password: "your-password",
//	}
//
// SendGrid:
//
//	Config{
//	    Host:     "smtp.sendgrid.net",
//	    Port:     587,
//	    Username: "apikey",
//	    Password: "your-sendgrid-api-key",
//	}
//
// Mailgun:
//
//	Config{
//	    Host:     "smtp.mailgun.org",
//	    Port:     587,
//	    Username: "postmaster@your-domain.mailgun.org",
//	    Password: "your-mailgun-password",
//	}
//
// # Configuration with Viper
//
// Example configuration file (configs/default.yaml):
//
//	smtp:
//	  host: smtp.gmail.com
//	  port: 587
//	  username: ${SMTP_USERNAME}
//	  password: ${SMTP_PASSWORD}
//	  from_email: noreply@example.com
//	  from_name: My Application
//
// Loading configuration:
//
//	import (
//	    "github.com/poly-workshop/go-webmods/app"
//	    "github.com/poly-workshop/go-webmods/smtp_mailer"
//	)
//
//	app.Init(".")
//	cfg := app.Config()
//
//	mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
//	    Host:      cfg.GetString("smtp.host"),
//	    Port:      cfg.GetInt("smtp.port"),
//	    Username:  cfg.GetString("smtp.username"),
//	    Password:  cfg.GetString("smtp.password"),
//	    FromEmail: cfg.GetString("smtp.from_email"),
//	    FromName:  cfg.GetString("smtp.from_name"),
//	})
//
// # Security
//
// The mailer uses TLS for secure communication:
//   - Connects via plain TCP initially
//   - Upgrades to TLS using STARTTLS
//   - Enforces minimum TLS version 1.2
//   - Validates server certificate
//
// # Email Templates
//
// For HTML emails with dynamic content, use Go's html/template:
//
//	import (
//	    "bytes"
//	    "html/template"
//	)
//
//	tmpl, err := template.ParseFiles("templates/welcome.html")
//	if err != nil {
//	    return err
//	}
//
//	var body bytes.Buffer
//	err = tmpl.Execute(&body, struct {
//	    Name string
//	    URL  string
//	}{
//	    Name: "Alice",
//	    URL:  "https://example.com/verify?token=abc123",
//	})
//	if err != nil {
//	    return err
//	}
//
//	err = mailer.SendEmail(smtp_mailer.Message{
//	    ToEmails: []string{"alice@example.com"},
//	    Subject:  "Verify your email",
//	    Body:     body.String(),
//	    IsHTML:   true,
//	})
//
// # Error Handling
//
// SendEmail returns an error if:
//   - Connection to SMTP server fails
//   - TLS upgrade fails
//   - Authentication fails
//   - Email is rejected by the server
//   - Network errors occur during sending
//
// Always check the error:
//
//	err := mailer.SendEmail(msg)
//	if err != nil {
//	    log.Printf("Failed to send email: %v", err)
//	    // Handle error (retry, log, notify admin, etc.)
//	    return err
//	}
//
// # Best Practices
//
//   - Store SMTP credentials in environment variables or secret management
//   - Use app passwords for Gmail (not your account password)
//   - Set appropriate From name for better deliverability
//   - Validate email addresses before sending
//   - Use HTML templates for consistent styling
//   - Include plain text alternatives for HTML emails (not currently supported, consider adding)
//   - Implement retry logic for transient failures
//   - Rate limit email sending to avoid being flagged as spam
//   - Use dedicated email service providers (SendGrid, Mailgun) for production
//   - Test emails thoroughly before production deployment
//
// # Rate Limiting
//
// Implement rate limiting for bulk emails:
//
//	import "time"
//
//	for _, recipient := range recipients {
//	    err := mailer.SendEmail(smtp_mailer.Message{
//	        ToEmails: []string{recipient},
//	        Subject:  "Newsletter",
//	        Body:     content,
//	    })
//	    if err != nil {
//	        log.Printf("Failed to send to %s: %v", recipient, err)
//	    }
//	    time.Sleep(100 * time.Millisecond)  // Rate limit
//	}
//
// # Limitations
//
// Current limitations:
//   - No attachment support
//   - No CC/BCC support
//   - No plain text alternative for HTML emails
//   - No email validation
//   - No retry mechanism
//   - No bulk sending optimization
//
// For advanced features, consider using a dedicated email library or service.
//
// # Example: Password Reset Email
//
//	func sendPasswordReset(mailer *smtp_mailer.Mailer, email, token string) error {
//	    resetURL := fmt.Sprintf("https://example.com/reset?token=%s", token)
//
//	    body := fmt.Sprintf(`
//	        <html>
//	            <body style="font-family: Arial, sans-serif;">
//	                <h2>Password Reset Request</h2>
//	                <p>Click the link below to reset your password:</p>
//	                <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px;">
//	                    Reset Password
//	                </a>
//	                <p>This link will expire in 1 hour.</p>
//	                <p>If you didn't request this, please ignore this email.</p>
//	            </body>
//	        </html>
//	    `, resetURL)
//
//	    return mailer.SendEmail(smtp_mailer.Message{
//	        ToEmails: []string{email},
//	        Subject:  "Password Reset Request",
//	        Body:     body,
//	        IsHTML:   true,
//	    })
//	}
package smtp_mailer
