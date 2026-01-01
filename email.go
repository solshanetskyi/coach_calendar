package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"
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

// generateICalendar creates an iCalendar (ICS) format string for the appointment
func generateICalendar(name, email string, slotTime time.Time) string {
	// Calculate end time (30 minutes after start)
	endTime := slotTime.Add(30 * time.Minute)

	// Format times in iCalendar format (YYYYMMDDTHHMMSSZ in UTC)
	startUTC := slotTime.UTC().Format("20060102T150405Z")
	endUTC := endTime.UTC().Format("20060102T150405Z")
	now := time.Now().UTC().Format("20060102T150405Z")

	// Generate a unique ID for the event
	eventID := fmt.Sprintf("%d@coach-calendar.com", time.Now().UnixNano())

	// Create iCalendar content
	ical := fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
PRODID:-//Coach Calendar//Booking System//EN
CALSCALE:GREGORIAN
METHOD:REQUEST
BEGIN:VEVENT
UID:%s
DTSTAMP:%s
DTSTART:%s
DTEND:%s
SUMMARY:Coaching Session with %s
DESCRIPTION:Your coaching appointment has been confirmed.\n\nClient: %s\nEmail: %s
LOCATION:Online/TBD
STATUS:CONFIRMED
SEQUENCE:0
BEGIN:VALARM
TRIGGER:-PT15M
ACTION:DISPLAY
DESCRIPTION:Reminder: Coaching session in 15 minutes
END:VALARM
END:VEVENT
END:VCALENDAR`, eventID, now, startUTC, endUTC, name, name, email)

	return ical
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

A calendar invitation is attached to this email. You can add this appointment to your calendar by opening the attachment.

Please make sure to arrive on time for your appointment.

If you need to cancel or reschedule, please contact us as soon as possible.

Best regards,
Coach Calendar Team

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
This is an automated message. Please do not reply to this email.
`, name, formattedTime, name, email)

	// Generate iCalendar attachment
	icalContent := generateICalendar(name, email, slotTime)

	// Send via AWS SES or SMTP
	if e.UseAWSSES {
		return e.sendViaSES(email, subject, body, icalContent)
	}
	return e.sendViaSMTP(email, subject, body, icalContent)
}

func (e *EmailService) sendViaSES(toEmail, subject, body, icalContent string) error {
	// Create boundary for multipart message
	boundary := fmt.Sprintf("boundary_%d", rand.Int63())

	// Build raw email message with attachment
	rawMessage := fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="%s"

--%s
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 7bit

%s

--%s
Content-Type: text/calendar; charset=UTF-8; method=REQUEST; name="invite.ics"
Content-Transfer-Encoding: base64
Content-Disposition: attachment; filename="invite.ics"

%s
--%s--`, e.From, toEmail, subject, boundary, boundary, body, boundary, base64.StdEncoding.EncodeToString([]byte(icalContent)), boundary)

	// Send raw email via SES
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: []byte(rawMessage),
		},
	}

	result, err := e.SESClient.SendRawEmail(input)
	if err != nil {
		log.Printf("Failed to send confirmation email via AWS SES to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send confirmation email via AWS SES: %w", err)
	}

	log.Printf("Confirmation email sent successfully via AWS SES to %s (MessageId: %s)", toEmail, *result.MessageId)
	return nil
}

func (e *EmailService) sendViaSMTP(toEmail, subject, body, icalContent string) error {
	// Create boundary for multipart message
	boundary := fmt.Sprintf("boundary_%d", rand.Int63())

	// Build multipart email with calendar attachment
	var message strings.Builder
	message.WriteString(fmt.Sprintf("From: %s\r\n", e.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", toEmail))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	message.WriteString("\r\n")

	// Text body part
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)
	message.WriteString("\r\n\r\n")

	// Calendar attachment part
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/calendar; charset=UTF-8; method=REQUEST; name=\"invite.ics\"\r\n")
	message.WriteString("Content-Transfer-Encoding: base64\r\n")
	message.WriteString("Content-Disposition: attachment; filename=\"invite.ics\"\r\n")
	message.WriteString("\r\n")
	message.WriteString(base64.StdEncoding.EncodeToString([]byte(icalContent)))
	message.WriteString("\r\n")
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	// Set up authentication
	auth := smtp.PlainAuth("", e.From, e.Password, e.SMTPHost)

	// Send the email
	addr := fmt.Sprintf("%s:%s", e.SMTPHost, e.SMTPPort)
	err := smtp.SendMail(addr, auth, e.From, []string{toEmail}, []byte(message.String()))
	if err != nil {
		log.Printf("Failed to send confirmation email via SMTP to %s: %v", toEmail, err)
		return fmt.Errorf("failed to send confirmation email via SMTP: %w", err)
	}

	log.Printf("Confirmation email sent successfully via SMTP to %s", toEmail)
	return nil
}
