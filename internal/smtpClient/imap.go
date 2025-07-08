package smtpclient

import (
	"fmt"
	"sync"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// IMAPConfig holds the configuration for the IMAP client
type IMAPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// IMAPClient represents an IMAP client that can connect to a mail server
type IMAPClient struct {
	config IMAPConfig
	client *client.Client
	mu     sync.Mutex
}

// NewIMAPClient creates a new IMAP client with the given configuration
func NewIMAPClient(config IMAPConfig) *IMAPClient {
	return &IMAPClient{
		config: config,
	}
}

// Connect establishes a connection to the IMAP server
func (c *IMAPClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Connect to server without TLS for local development/testing
	imapClient, err := client.Dial(fmt.Sprintf("%s:%d", c.config.Host, c.config.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP server: %w", err)
	}

	// Login
	if err := imapClient.Login(c.config.Username, c.config.Password); err != nil {
		imapClient.Logout()
		return fmt.Errorf("failed to login to IMAP server: %w", err)
	}

	c.client = imapClient
	return nil
}

// Disconnect closes the connection to the IMAP server
func (c *IMAPClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil
	}

	if err := c.client.Logout(); err != nil {
		return fmt.Errorf("failed to logout from IMAP server: %w", err)
	}

	c.client = nil
	return nil
}

// GetInbox retrieves messages from the user's inbox
func (c *IMAPClient) GetInbox(limit int) ([]Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil, fmt.Errorf("not connected to IMAP server")
	}

	// Select INBOX
	mbox, err := c.client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %w", err)
	}

	// Get the last 'limit' messages
	from := uint32(1)
	to := mbox.Messages
	if to > uint32(limit) && to > 0 {
		from = to - uint32(limit) + 1
	}

	if from > to {
		return []Message{}, nil
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	// Get message envelope and flags
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	var result []Message
	for msg := range messages {
		message := Message{
			ID:      fmt.Sprintf("%d", msg.SeqNum),
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date,
		}

		// Set From
		if len(msg.Envelope.From) > 0 {
			message.From = msg.Envelope.From[0].Address()
		}

		// Set To
		for _, addr := range msg.Envelope.To {
			message.To = append(message.To, addr.Address())
		}

		result = append(result, message)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %w", err)
	}

	return result, nil
}
