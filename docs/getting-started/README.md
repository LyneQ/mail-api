# Getting Started with MailAPI

This guide will help you set up and run the MailAPI project on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.16 or later)
- Git
- SQLite

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/mailapi.git
cd mailapi
```

2. Install dependencies:

```bash
go mod download
```

## Configuration

1. Copy the example configuration file:

```bash
cp config/config.example.ini config/config.ini
```

2. Edit the configuration file to match your environment:

```ini
[AllowedDomains]
domains = localhost, yourdomain.com

[Database]
driver = sqlite3
path = ./db/dev.db

[Api]
port = 1323

[SMTP]
host = smtp.example.com
port = 587
username = your_username
password = your_password

[IMAP]
host = imap.example.com
port = 993
username = your_username
password = your_password
```

## Running the Application

To start the MailAPI server:

```bash
go run main.go
```

The API will be available at `http://localhost:1323`.

## Testing

To run the tests:

```bash
go test ./test/...
```

## API Endpoints

Once the server is running, you can access the following endpoints:

- `GET /` - Hello World (test endpoint)
- `POST /api/signup` - Register a new user
- `POST /api/signin` - Login a user
- `GET /api/me` - Get current user info (requires authentication)
- `GET /api/signout` - Logout (requires authentication)
- `GET /api/email/inbox` - Get inbox emails (requires authentication)
- `GET /api/email/:id` - Get a specific email (requires authentication)
- `POST /api/email/send` - Send an email (requires authentication)

For more detailed API documentation, see the [API Reference](../api/README.md).

## Next Steps

- Learn about the [API Reference](../api/README.md)
- Explore the [Configuration Options](../config/README.md)
- Understand the [Architecture](../architecture/README.md)