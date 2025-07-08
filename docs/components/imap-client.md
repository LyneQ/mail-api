# IMAP Client Component

The IMAP (Internet Message Access Protocol) client component is responsible for retrieving emails from mail servers in the MailAPI application.

## Overview

The IMAP client provides functionality for connecting to an IMAP server, retrieving emails from the inbox, and fetching specific emails with their full content and attachments.

## Implementation

The IMAP client is implemented in the `internal/smtpClient/imap.go` file. It uses the `github.com/emersion/go-imap` package to handle the low-level IMAP protocol details.

### Key Structures

#### IMAPConfig

```
type IMAPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}
```

This structure holds the configuration for connecting to an IMAP server.

#### IMAPClient

```
type IMAPClient struct {
    config IMAPConfig
    client *client.Client
    mu     sync.Mutex
}
```

The IMAPClient structure encapsulates the IMAP client functionality. It includes:
- `config`: The IMAP server configuration
- `client`: The go-imap client used to connect to the IMAP server
- `mu`: A mutex for thread safety

#### Message

The IMAP client uses the same Message structure as the SMTP client:

```
type Message struct {
    ID          string
    From        string
    To          []string
    Subject     string
    Body        string
    Date        time.Time
    Attachments []Attachment
}
```

### Key Methods

#### NewIMAPClient

```
func NewIMAPClient(config IMAPConfig) *IMAPClient
```

Creates a new IMAP client with the given configuration.

#### Connect

```
func (c *IMAPClient) Connect() error
```

Establishes a connection to the IMAP server and logs in with the configured credentials. This method is thread-safe.

#### Disconnect

```
func (c *IMAPClient) Disconnect() error
```

Closes the connection to the IMAP server by logging out. This method is thread-safe.

#### GetInbox

```
func (c *IMAPClient) GetInbox(limit int) ([]Message, error)
```

Retrieves a specified number of messages from the user's inbox. This method:
1. Selects the INBOX mailbox
2. Fetches the most recent messages (up to the specified limit)
3. Extracts basic information like sender, recipients, subject, and date
4. Returns the messages as a slice of Message structures

#### GetEmailByID

```
func (c *IMAPClient) GetEmailByID(id string) (*Message, error)
```

Retrieves a specific email by its ID with full details. This method:
1. Selects the INBOX mailbox
2. Fetches the specified message with its full content
3. Parses the message body and any attachments
4. Returns a Message structure with all the email details

## Configuration

The IMAP client is configured through the `config/config.ini` file in the `[IMAP]` section:

```
[IMAP]
host = imap.example.com
port = 993
username = your_username
password = your_password
```

## Usage Example

Here's an example of how the IMAP client is used in the application:

```
// Create a new IMAP client
config := smtpclient.IMAPConfig{
    Host:     "imap.example.com",
    Port:     993,
    Username: "user@example.com",
    Password: "password",
}
client := smtpclient.NewIMAPClient(config)

// Connect to the IMAP server
if err := client.Connect(); err != nil {
    log.Fatalf("Failed to connect to IMAP server: %v", err)
}
defer client.Disconnect()

// Get the most recent 50 emails from the inbox
messages, err := client.GetInbox(50)
if err != nil {
    log.Fatalf("Failed to get inbox: %v", err)
}

// Print the subject of each email
for _, msg := range messages {
    fmt.Printf("Email: %s - From: %s\n", msg.Subject, msg.From)
}

// Get a specific email by ID
email, err := client.GetEmailByID("1")
if err != nil {
    log.Fatalf("Failed to get email: %v", err)
}
fmt.Printf("Email body: %s\n", email.Body)
```

## Security Considerations

The IMAP client includes several security features:

1. **TLS Encryption**: The client typically connects to IMAP servers on port 993, which uses implicit TLS.
2. **Authentication**: The client authenticates with the IMAP server using the provided username and password.
3. **Thread Safety**: The client uses a mutex to ensure thread safety when multiple goroutines are accessing the IMAP server.

## Limitations

The current implementation has a few limitations:

1. **No Connection Pooling**: Each client maintains a single connection to the IMAP server.
2. **Limited Mailbox Support**: The client only works with the INBOX mailbox, not other folders.
3. **No Idle Support**: The client doesn't support the IMAP IDLE command for real-time notifications.
4. **No Caching**: Emails are fetched from the server each time, with no local caching.

These limitations could be addressed in future versions of the component.

## Error Handling

The IMAP client includes comprehensive error handling:

1. Connection errors are returned when the client fails to connect to the server
2. Authentication errors are returned when login fails
3. Mailbox selection errors are returned when the INBOX can't be selected
4. Fetch errors are returned when messages can't be retrieved
5. Parse errors are returned when message content can't be parsed

All errors include descriptive messages to help with debugging.