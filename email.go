package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type EmailService struct {
	SMTPHost  string
	SMTPPort  string
	From      string
	Password  string
	UseAWSSES bool
	SESClient *ses.SES
	AWSRegion string
	Enabled   bool
}

func NewEmailService() *EmailService {
	// Check if AWS SES should be used
	useAWSSES := os.Getenv("USE_AWS_SES") == "true"
	from := os.Getenv("SMTP_FROM") // Used for both SMTP and SES

	var sesClient *ses.SES
	var awsRegion string
	var enabled bool

	if useAWSSES {
		// AWS SES Configuration
		awsRegion = os.Getenv("AWS_REGION")
		if awsRegion == "" {
			awsRegion = "us-east-1" // Default region
		}

		if from == "" {
			log.Println("Email service disabled - SMTP_FROM not configured")
			log.Println("To enable AWS SES email confirmations, set: USE_AWS_SES=true, SMTP_FROM, AWS_REGION")
			enabled = false
		} else {
			// Create AWS session
			sess, err := session.NewSession(&aws.Config{
				Region: aws.String(awsRegion),
			})
			if err != nil {
				log.Printf("Failed to create AWS session: %v", err)
				log.Println("Email service disabled - AWS session creation failed")
				enabled = false
			} else {
				sesClient = ses.New(sess)
				enabled = true
				log.Printf("Email service enabled using AWS SES in region: %s", awsRegion)
			}
		}
	} else {
		// SMTP Configuration
		host := os.Getenv("SMTP_HOST")
		port := os.Getenv("SMTP_PORT")
		password := os.Getenv("SMTP_PASSWORD")

		enabled = host != "" && port != "" && from != "" && password != ""

		if !enabled {
			log.Println("Email service disabled - SMTP configuration not found")
			log.Println("To enable SMTP email confirmations, set: SMTP_HOST, SMTP_PORT, SMTP_FROM, SMTP_PASSWORD")
			log.Println("To enable AWS SES email confirmations, set: USE_AWS_SES=true, SMTP_FROM, AWS_REGION")
		} else {
			log.Println("Email service enabled using SMTP")
		}

		return &EmailService{
			SMTPHost: host,
			SMTPPort: port,
			From:     from,
			Password: password,
			Enabled:  enabled,
		}
	}

	return &EmailService{
		From:      from,
		UseAWSSES: useAWSSES,
		SESClient: sesClient,
		AWSRegion: awsRegion,
		Enabled:   enabled,
	}
}

func (e *EmailService) SendBookingConfirmation(name, email string, slotTime time.Time) error {
	if !e.Enabled {
		log.Printf("Email service disabled - skipping confirmation email to %s", email)
		return nil
	}

	// Format the booking time
	formattedTime := slotTime.Format("Monday, January 2, 2006 at 3:04 PM MST")

	// Create email subject and body
	subject := "Booking Confirmation - Coach Calendar"
	body := fmt.Sprintf(`Hello %s,

Thank you for booking an appointment!

Appointment Details:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📅 Date & Time: %s
⏱️  Duration: 30 minutes
👤 Name: %s
📧 Email: %s

Please make sure to arrive on time for your appointment.

If you need to cancel or reschedule, please contact us as soon as possible.

Best regards,
Coach Calendar Team

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
This is an automated message. Please do not reply to this email.
`, name, formattedTime, name, email)

	// Send via AWS SES or SMTP
	if e.UseAWSSES {
		return e.sendViaSES(email, subject, body)
	}
	return e.sendViaSMTP(email, subject, body)
}

func (e *EmailService) sendViaSES(toEmail, subject, body string) error {
	// Create SES input
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(toEmail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(e.From),
	}

	// Send the email via SES
	result, err := e.SESClient.SendEmail(input)
	if err != nil {
		log.Printf("Failed to send confirmation email via AWS SES to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send confirmation email via AWS SES: %w", err)
	}

	log.Printf("Confirmation email sent successfully via AWS SES to %s (MessageId: %s)", toEmail, *result.MessageId)
	return nil
}

func (e *EmailService) sendViaSMTP(toEmail, subject, body string) error {
	// Compose the email message
	message := fmt.Sprintf("From: %s\r\n", e.From)
	message += fmt.Sprintf("To: %s\r\n", toEmail)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// Set up authentication
	auth := smtp.PlainAuth("", e.From, e.Password, e.SMTPHost)

	// Send the email
	addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
	err := smtp.SendMail(addr, auth, e.From, []string{toEmail}, []byte(message))
	if err != nil {
		log.Printf("Failed to send confirmation email via SMTP to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send confirmation email via SMTP: %w", err)
	}

	log.Printf("Confirmation email sent successfully via SMTP to %s", toEmail)
	return nil
}
