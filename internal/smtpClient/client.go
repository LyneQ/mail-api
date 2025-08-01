package smtpclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

// SMTPConfig holds the configuration for the SMTP client
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// Client represents an SMTP client that can connect to a mail server
type Client struct {
	config SMTPConfig
	dialer *gomail.Dialer
	mu     sync.Mutex
}

// Message represents an email message
type Message struct {
	ID          string
	From        string
	To          []string
	Subject     string
	Body        string
	Date        time.Time
	Attachments []Attachment
	Flags       []string
	Size        uint32 // Size of the message in bytes
}

// Attachment represents an email attachment
type Attachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// NewClient creates a new SMTP client with the given configuration
func NewClient(config SMTPConfig) *Client {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)

	if config.Port == 1025 {
		dialer.SSL = false

		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	} else {
		// Use system's certificate pool for proton-bridge certificates
		dialer.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &Client{
		config: config,
		dialer: dialer,
	}
}

// Connect establishes a connection to the SMTP server
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, err := c.dialer.Dial()
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	return nil
}

// SendMessage sends an email message
func (c *Client) SendMessage(from string, to []string, subject, body string, attachments []Attachment) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	for _, attachment := range attachments {
		m.Attach(attachment.Filename,
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(attachment.Content)
				return err
			}),
			gomail.SetHeader(map[string][]string{
				"Content-Type": {attachment.MimeType},
			}),
		)
	}

	if err := c.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
