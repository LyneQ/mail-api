package smtpclient

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
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

	var imapClient *client.Client
	var err error

	// For port 1143, use non-TLS connection with STARTTLS
	if c.config.Port == 1143 {

		imapClient, err = client.Dial(fmt.Sprintf("%s:%d", c.config.Host, c.config.Port))
		if err != nil {
			return fmt.Errorf("failed to connect to IMAP server: %w", err)
		}

		certPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile("config/tls/cert.pem")
		if err != nil {
			return fmt.Errorf("failed to read certificate: %w", err)
		}
		certPool.AppendCertsFromPEM(cert)

		tlsConfig := &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
		}

		if err := imapClient.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	} else {
		// For other ports, use TLS
		// Load TLS certificate
		certPool := x509.NewCertPool()
		cert, err := ioutil.ReadFile("config/tls/cert.pem")
		if err != nil {
			return fmt.Errorf("failed to read certificate: %w", err)
		}
		certPool.AppendCertsFromPEM(cert)

		tlsConfig := &tls.Config{
			RootCAs:            certPool,
			InsecureSkipVerify: false,
		}

		imapClient, err = client.DialTLS(fmt.Sprintf("%s:%d", c.config.Host, c.config.Port), tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to IMAP server: %w", err)
		}
	}

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

	mailboxes := make(chan *imap.MailboxInfo, 50)
	done := make(chan error, 1)

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

	mbox, err := c.client.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %w", err)
	}

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

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	fmt.Println("Fetching inbox messages")

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

	mbox, err := c.client.Select(folder, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select folder: %w", err)
	}

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

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags}
	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	fmt.Println("Fetching messages from folder:", folder)

	var result []Message
	for msg := range messages {

		message := Message{
			ID:      fmt.Sprintf("%d", msg.SeqNum),
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date,
			Flags:   msg.Flags,
		}

		if len(msg.Envelope.From) > 0 {
			message.From = msg.Envelope.From[0].Address()
		}

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
func (c *IMAPClient) GetEmailByID(id string, folder string) (*Message, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client == nil {
		return nil, fmt.Errorf("not connected to IMAP server")
	}

	seqNum, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid email ID: %w", err)
	}

	if folder == "" {
		folder = "INBOX"
	}

	fmt.Println("Selecting folder for email ID", id, ":", folder)
	_, err = c.client.Select(folder, false)
	if err != nil {
		return nil, fmt.Errorf("failed to select folder %s: %w", folder, err)
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uint32(seqNum))

	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchBodyStructure, "BODY[]"}
	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- c.client.Fetch(seqSet, items, messages)
	}()

	var message *Message
	for msg := range messages {
		message = &Message{
			ID:      fmt.Sprintf("%d", msg.SeqNum),
			Subject: msg.Envelope.Subject,
			Date:    msg.Envelope.Date,
			Flags:   msg.Flags,
		}

		if len(msg.Envelope.From) > 0 {
			message.From = msg.Envelope.From[0].Address()
		}

		for _, addr := range msg.Envelope.To {
			message.To = append(message.To, addr.Address())
		}

		for _, literal := range msg.Body {
			mr, err := mail.CreateReader(literal)
			if err != nil {
				continue
			}

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
					b, _ := ioutil.ReadAll(p.Body)
					message.Body = string(b)
				case *mail.AttachmentHeader:
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
