package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strconv"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
}

func NewEmailService() *EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}

	return &EmailService{
		host:     os.Getenv("SMTP_HOST"),
		port:     port,
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     os.Getenv("SMTP_FROM"),
	}
}

func (es *EmailService) SendVerificationEmail(to, email, token string) error {
	subject := "Email Verification - Secure Auth System"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Email Verification</h2>
				<p>Hello,</p>
				<p>Please verify your email address by using the verification code below:</p>
				<h3 style="color: #007bff;">%s</h3>
				<p>This code will expire in 10 minutes.</p>
				<p>If you did not register for this account, please ignore this email.</p>
				<hr>
				<p>Secure Auth System</p>
			</body>
		</html>
	`, token)

	return es.sendEmail(to, subject, body)
}

func (es *EmailService) SendPhoneOTPEmail(to, otp string) error {
	subject := "Phone Verification OTP - Secure Auth System"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Phone Verification Code</h2>
				<p>Hello,</p>
				<p>Your phone verification code is:</p>
				<h3 style="color: #28a745; font-size: 32px; letter-spacing: 5px;">%s</h3>
				<p>This code will expire in 5 minutes.</p>
				<p>Do not share this code with anyone.</p>
				<p>If you did not request this code, please ignore this email.</p>
				<hr>
				<p>Secure Auth System</p>
			</body>
		</html>
	`, otp)

	return es.sendEmail(to, subject, body)
}

func (es *EmailService) SendLoginOTPEmail(to, otp string) error {
	subject := "Login Verification Code - Secure Auth System"
	body := fmt.Sprintf(`
		<html>
			<body>
				<h2>Login Verification Code</h2>
				<p>Hello,</p>
				<p>Your login verification code is:</p>
				<h3 style="color: #007bff; font-size: 32px; letter-spacing: 5px;">%s</h3>
				<p>This code will expire in 5 minutes.</p>
				<p>If you did not attempt to log in, please ignore this email.</p>
				<hr>
				<p>Secure Auth System</p>
			</body>
		</html>
	`, otp)

	return es.sendEmail(to, subject, body)
}

func (es *EmailService) sendEmail(to, subject, body string) error {
	if es.host == "" || es.username == "" || es.password == "" {
		log.Println("[MOCK EMAIL] To:", to, "Subject:", subject)
		return nil
	}

	auth := smtp.PlainAuth("", es.username, es.password, es.host)
	addr := fmt.Sprintf("%s:%d", es.host, es.port)

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nMIME-Version: 1.0\nContent-Type: text/html; charset=UTF-8\n\n%s",
		es.from, to, subject, body)

	err := smtp.SendMail(addr, auth, es.from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("Error sending email to %s: %v", to, err)
		return err
	}

	log.Printf("Email sent successfully to %s", to)
	return nil
}
