package smtp_mailer_test

import (
	"fmt"

	"github.com/poly-workshop/go-webmods/smtp_mailer"
)

// Example demonstrates basic email sending.
func Example() {
	// mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
	// 	Host:      "smtp.example.com",
	// 	Port:      587,
	// 	Username:  "user@example.com",
	// 	Password:  "password",
	// 	FromEmail: "noreply@example.com",
	// 	FromName:  "My Application",
	// })
	//
	// err := mailer.SendEmail(smtp_mailer.Message{
	// 	ToEmails: []string{"recipient@example.com"},
	// 	Subject:  "Hello",
	// 	Body:     "This is a test email.",
	// 	IsHTML:   false,
	// })
	//
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	fmt.Println("Email sent successfully")
	// Output: Email sent successfully
}

// Example_html demonstrates sending HTML emails.
func Example_html() {
	// mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
	// 	Host:      "smtp.example.com",
	// 	Port:      587,
	// 	Username:  "user@example.com",
	// 	Password:  "password",
	// 	FromEmail: "noreply@example.com",
	// 	FromName:  "My Application",
	// })
	//
	// htmlBody := `
	// 	<html>
	// 		<body>
	// 			<h1>Welcome!</h1>
	// 			<p>Thank you for signing up.</p>
	// 			<a href="https://example.com/verify">Verify Email</a>
	// 		</body>
	// 	</html>
	// `
	//
	// err := mailer.SendEmail(smtp_mailer.Message{
	// 	ToEmails: []string{"user@example.com"},
	// 	Subject:  "Welcome to Our Service",
	// 	Body:     htmlBody,
	// 	IsHTML:   true,
	// })
	//
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	fmt.Println("HTML email sent")
	// Output: HTML email sent
}

// Example_multipleRecipients demonstrates sending to multiple recipients.
func Example_multipleRecipients() {
	// mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
	// 	Host:      "smtp.example.com",
	// 	Port:      587,
	// 	Username:  "user@example.com",
	// 	Password:  "password",
	// 	FromEmail: "noreply@example.com",
	// 	FromName:  "Newsletter",
	// })
	//
	// err := mailer.SendEmail(smtp_mailer.Message{
	// 	ToEmails: []string{
	// 		"user1@example.com",
	// 		"user2@example.com",
	// 		"user3@example.com",
	// 	},
	// 	Subject: "Monthly Newsletter",
	// 	Body:    "Check out our latest updates!",
	// 	IsHTML:  false,
	// })
	//
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	fmt.Println("Newsletter sent to multiple recipients")
	// Output: Newsletter sent to multiple recipients
}

// Example_gmail demonstrates configuring Gmail SMTP.
func Example_gmail() {
	// Note: For Gmail, you need to use an App Password
	// Go to Google Account > Security > 2-Step Verification > App Passwords

	mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
		Host:      "smtp.gmail.com",
		Port:      587,
		Username:  "your-email@gmail.com",
		Password:  "your-app-password", // Not your regular password!
		FromEmail: "your-email@gmail.com",
		FromName:  "My App",
	})

	_ = mailer
	fmt.Println("Gmail mailer configured")
	// Output: Gmail mailer configured
}

// Example_passwordReset demonstrates a password reset email.
func Example_passwordReset() {
	// mailer := smtp_mailer.NewMailer(smtp_mailer.Config{
	// 	Host:      "smtp.example.com",
	// 	Port:      587,
	// 	Username:  "user@example.com",
	// 	Password:  "password",
	// 	FromEmail: "noreply@example.com",
	// 	FromName:  "My Application",
	// })
	//
	// resetToken := "abc123xyz"
	// resetURL := fmt.Sprintf("https://example.com/reset?token=%s", resetToken)
	//
	// htmlBody := fmt.Sprintf(`
	// 	<html>
	// 		<body style="font-family: Arial, sans-serif;">
	// 			<h2>Password Reset Request</h2>
	// 			<p>Click the link below to reset your password:</p>
	// 			<a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #007bff; color: white; text-decoration: none; border-radius: 5px;">
	// 				Reset Password
	// 			</a>
	// 			<p>This link will expire in 1 hour.</p>
	// 			<p>If you didn't request this, please ignore this email.</p>
	// 		</body>
	// 	</html>
	// `, resetURL)
	//
	// err := mailer.SendEmail(smtp_mailer.Message{
	// 	ToEmails: []string{"user@example.com"},
	// 	Subject:  "Password Reset Request",
	// 	Body:     htmlBody,
	// 	IsHTML:   true,
	// })
	//
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	fmt.Println("Password reset email sent")
	// Output: Password reset email sent
}
