package smtpclient

import (
	"fmt"
	"io"
	"io/ioutil" // Using ioutil for simplicity, though it's deprecated in favor of io and os packages in newer Go versions
	"strconv"
	"sync"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
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

// GetFolders retourne la liste des boîtes aux lettres (mailboxes) disponibles sur le serveur IMAP.
func (c *IMAPClient) GetFolders() ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil, fmt.Errorf("not connected to IMAP server")
	}

	// Canal pour récupérer les résultats
	mailboxes := make(chan *imap.MailboxInfo, 50)
	done := make(chan error, 1)

	// Lance la requête IMAP pour lister les boîtes.
	// Ex: "" pour le "reference", "*" pour le "mailbox pattern"
	go func() {
		done <- c.client.List("", "*", mailboxes)
	}()

	var folderNames []string
	for m := range mailboxes {
		folderNames = append(folderNames, m.Name)
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to list mailboxes: %w", err)
	}

	return folderNames, nil
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
			Flags:   msg.Flags,
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

func (c *IMAPClient) GetFolderMessages(folder string, limit int) ([]Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil, fmt.Errorf("not connected to IMAP server")
	}

	// Select INBOX
	mbox, err := c.client.Select(folder, false)
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
			Flags:   msg.Flags,
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

// GetEmailByID retrieves a specific email by its ID with full details
func (c *IMAPClient) GetEmailByID(id string) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil, fmt.Errorf("not connected to IMAP server")
	}

	// Convert ID to sequence number
	seqNum, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid email ID: %w", err)
	}

	// Select INBOX
	_, err = c.client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %w", err)
	}

	// Create sequence set for this specific message
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uint32(seqNum))

	// Get message envelope, flags, and body structure
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY[]"}
	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	var message *Message
	for msg := range messages {
		// Create basic message with envelope data
		message = &Message{
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

		// Get the body
		for _, literal := range msg.Body {
			// Parse the message
			mr, err := mail.CreateReader(literal)
			if err != nil {
				continue
			}

			// Process each part
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				if err != nil {
					continue
				}

				switch h := p.Header.(type) {
				case *mail.InlineHeader:
					// This is the message body
					b, _ := ioutil.ReadAll(p.Body)
					message.Body = string(b)
				case *mail.AttachmentHeader:
					// This is an attachment
					filename, _ := h.Filename()
					b, _ := ioutil.ReadAll(p.Body)
					contentType, _, _ := h.ContentType()

					message.Attachments = append(message.Attachments, Attachment{
						Filename: filename,
						Content:  b,
						MimeType: contentType,
					})
				}
			}
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch message: %w", err)
	}

	if message == nil {
		return nil, fmt.Errorf("message with ID %s not found", id)
	}

	return message, nil
}
