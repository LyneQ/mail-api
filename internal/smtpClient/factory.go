package smtpclient

import (
	"strconv"

	"github.com/lyneq/mailapi/config"
)

// NewSMTPClientFromConfig creates a new SMTP client using the application configuration
func NewSMTPClientFromConfig() *Client {
	smtpConfig := config.GetSMTPConfig()
	port, err := strconv.Atoi(smtpConfig.Port)
	if err != nil {
		port = 25 // Default SMTP port
	}

	return NewClient(SMTPConfig{
		Host:     smtpConfig.Host,
		Port:     port,
		Username: smtpConfig.Username,
		Password: smtpConfig.Password,
	})
}

// NewIMAPClientFromConfig creates a new IMAP client using the application configuration
func NewIMAPClientFromConfig() *IMAPClient {
	imapConfig := config.GetIMAPConfig()
	port, err := strconv.Atoi(imapConfig.Port)
	if err != nil {
		port = 143 // Default IMAP port
	}

	return NewIMAPClient(IMAPConfig{
		Host:     imapConfig.Host,
		Port:     port,
		Username: imapConfig.Username,
		Password: imapConfig.Password,
	})
}
