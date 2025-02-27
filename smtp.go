package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"strings"

	gomail "github.com/wneessen/go-mail"
)

func sendEmail(smtpConfig *SmtpServerConfig, sender string, recipients []string, subject string, body string) error {
	email := gomail.NewMsg()
	// Populate email details
	if err := email.From(sender); err != nil {
		return fmt.Errorf("failed to add from to email: %v", err)
	}

	if err := email.ReplyTo(sender); err != nil {
		return fmt.Errorf("failed to add reply-to to email: %v", err)
	}

	if err := email.To(recipients...); err != nil {
		return fmt.Errorf("failed to add recipients to email: %v", err)
	}

	email.Subject(subject)
	email.SetBodyString(gomail.TypeTextPlain, body)

	// Set SMTP details
	clientOptions := []gomail.Option{
		gomail.WithPort(smtpConfig.Port),
	}

	if smtpConfig.Helo != "" {
		clientOptions = append(clientOptions, gomail.WithHELO(smtpConfig.Helo))
	}

	if smtpConfig.Auth {
		clientOptions = append(
			clientOptions,
			gomail.WithSMTPAuth(gomail.SMTPAuthLogin),
			gomail.WithUsername(smtpConfig.User),
			gomail.WithPassword(smtpConfig.Pass),
		)
	}

	switch smtpConfig.Tls {
	case ConstTlsNone:
		clientOptions = append(clientOptions, gomail.WithTLSPolicy(gomail.NoTLS))
	case ConstTlsInsecure:
		tlsSkipConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         smtpConfig.Server,
		}
		clientOptions = append(clientOptions, gomail.WithTLSConfig(tlsSkipConfig))
	case ConstTls:
		clientOptions = append(clientOptions, gomail.WithTLSPolicy(gomail.TLSMandatory))
	}

	// Create a new client using the options
	client, err := gomail.NewClient(smtpConfig.Server, clientOptions...)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}

	if client == nil {
		return fmt.Errorf("SMTP client is nil")
	}

	if err := client.DialAndSend(email); err != nil {
		return fmt.Errorf("failed to dial and send: %v", err)
	}

	log.Printf("Email sent successfully from %s to %s\n", sender, strings.Join(recipients, ", "))
	return nil
}
