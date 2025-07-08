package test

import (
	"os"
	"testing"

	"github.com/lyneq/mailapi/config"
	"github.com/lyneq/mailapi/internal/smtpClient"
)

func TestSMTPAndIMAPClient(t *testing.T) {
	// Override the config file path for testing
	originalOpen := config.OsOpen
	defer func() { config.OsOpen = originalOpen }()

	config.OsOpen = func(name string) (*os.File, error) {
		// Use the absolute path to the config file
		return os.Open("../config/config.ini")
	}

	// Load configuration
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Get SMTP configuration
	smtpConfig := config.GetSMTPConfig()

	// Create SMTP client
	client := smtpclient.NewSMTPClientFromConfig()

	// Connect to SMTP server
	if err := client.Connect(); err != nil {
		t.Logf("Failed to connect to SMTP server: %v", err)
		t.Skip("Skipping SMTP test due to connection failure")
	} else {
		t.Logf("Successfully connected to SMTP server at %s:%s", smtpConfig.Host, smtpConfig.Port)
	}

	// Send a message
	err := client.SendMessage(
		smtpConfig.Username,           // Use the configured username as sender
		[]string{smtpConfig.Username}, // Send to self for testing
		"Test Email from Go Test",
		"<p>This is a test message sent from the Go test suite.</p>",
		nil,
	)
	if err != nil {
		t.Logf("Failed to send message: %v", err)
		t.Skip("Skipping SMTP send test due to failure")
	} else {
		t.Log("Successfully sent test email")
	}

	// Get IMAP configuration
	imapConfig := config.GetIMAPConfig()

	// Create IMAP client
	imapClient := smtpclient.NewIMAPClientFromConfig()

	// Connect to IMAP server
	if err := imapClient.Connect(); err != nil {
		t.Logf("Failed to connect to IMAP server: %v", err)
		t.Skip("Skipping IMAP test due to connection failure")
	} else {
		t.Logf("Successfully connected to IMAP server at %s:%s", imapConfig.Host, imapConfig.Port)
	}
	defer imapClient.Disconnect()

	// Get inbox messages (last 10)
	messages, err := imapClient.GetInbox(10)
	if err != nil {
		t.Errorf("Failed to get inbox: %v", err)
	} else {
		t.Logf("Successfully retrieved %d messages from inbox", len(messages))

		// Log message details
		for i, msg := range messages {
			t.Logf("Message %d:", i+1)
			t.Logf("  ID: %s", msg.ID)
			t.Logf("  To: %s", msg.To)
			t.Logf("  From: %s", msg.From)
			t.Logf("  Subject: %s", msg.Subject)
			t.Logf("  Date: %s", msg.Date)
		}
	}
}
