# SMTP Client Component

The SMTP (Simple Mail Transfer Protocol) client component is responsible for sending emails from the MailAPI application.

## Overview

The SMTP client provides a clean interface for sending emails with or without attachments. It handles the connection to the SMTP server, authentication, and the actual sending of email messages.

## Implementation

The SMTP client is implemented in the `internal/smtpClient/client.go` file. It uses the `gopkg.in/gomail.v2` package to handle the low-level SMTP protocol details.

### Key Structures

#### SMTPConfig

```go
type SMTPConfig struct {
    Host     string
    Port     int
    Username string
    Password string
}
```

This structure holds the configuration for connecting to an SMTP server.

#### Client

```go
type Client struct {
    config SMTPConfig
    dialer *gomail.Dialer
    mu     sync.Mutex
}
```

The Client structure encapsulates the SMTP client functionality. It includes:
- `config`: The SMTP server configuration
- `dialer`: The gomail dialer used to connect to the SMTP server
- `mu`: A mutex for thread safety

#### Message

```go
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

This structure represents an email message with all its components.

#### Attachment

```go
type Attachment struct {
    Filename string
    Content  []byte
    MimeType string
}
```

This structure represents an email attachment with its filename, content, and MIME type.

### Key Methods

#### NewClient

```go
func NewClient(config SMTPConfig) *Client
```

Creates a new SMTP client with the given configuration.

#### Connect

```go
func (c *Client) Connect() error
```

Establishes a connection to the SMTP server. This method is thread-safe.

#### SendMessage

```go
func (c *Client) SendMessage(from string, to []string, subject, body string, attachments []Attachment) error
```

Sends an email message with the specified parameters. This method:
1. Creates a new gomail message
2. Sets the From, To, Subject, and Body fields
3. Adds any attachments
4. Connects to the SMTP server and sends the message

## Configuration

The SMTP client is configured through the `config/config.ini` file in the `[SMTP]` section:

```ini
[SMTP]
host = smtp.example.com
port = 587
username = your_username
password = your_password
```

## Usage Example

Here's an example of how the SMTP client is used in the application:

```go
// Create a new SMTP client
config := smtpclient.SMTPConfig{
    Host:     "smtp.example.com",
    Port:     587,
    Username: "user@example.com",
    Password: "password",
}
client := smtpclient.NewClient(config)

// Connect to the SMTP server
if err := client.Connect(); err != nil {
    log.Fatalf("Failed to connect to SMTP server: %v", err)
}

// Send an email
err := client.SendMessage(
    "sender@example.com",
    []string{"recipient@example.com"},
    "Hello",
    "<p>This is the email body</p>",
    []smtpclient.Attachment{
        {
            Filename: "document.pdf",
            Content:  []byte{...},
            MimeType: "application/pdf",
        },
    },
)
if err != nil {
    log.Fatalf("Failed to send email: %v", err)
}
```

## Security Considerations

The SMTP client includes several security features:

1. **TLS Encryption**: The client uses TLS to encrypt the connection to the SMTP server.
2. **Authentication**: The client authenticates with the SMTP server using the provided username and password.
3. **Thread Safety**: The client uses a mutex to ensure thread safety when multiple goroutines are sending emails.

Note that in the current implementation, TLS verification is skipped for local development/testing. In a production environment, this should be enabled.

## Limitations

The current implementation has a few limitations:

1. It doesn't support connection pooling, which could improve performance when sending many emails.
2. It doesn't implement retry logic for failed email sends.
3. It doesn't support DKIM or SPF for email authentication.

These limitations could be addressed in future versions of the component.