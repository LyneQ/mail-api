# MailAPI

MailAPI is a RESTful API service that provides email functionality through a simple HTTP interface. It allows applications to send and receive emails through SMTP and IMAP protocols, abstracting away the complexity of directly interacting with email servers.

## Key Features

- **Authentication**: Secure user authentication system
- **Email Operations**: Send emails and retrieve inbox messages
- **RESTful API**: Simple and consistent API design
- **Configurable**: (Coming soon) only `Proton Mail Bridge` is available for now
- **Multiple account**: (Coming soon) only one mail account is available for now

## Technology Stack

- **Backend**: Go (Golang)
- **Web Framework**: Echo
- **Database**: SQLite with GORM
- **Email Protocols**: SMTP (for sending) and IMAP (for receiving)

## Prerequisites

Before you begin, ensure you have the following installed:

- Go (version 1.24.4 or later)
- Git
- SQLite

## API Endpoints

### Authentication Endpoints

- `POST /api/signup` - Register a new user
- `POST /api/signin` - Login a user
- `GET /api/me` - Get current user info (requires authentication)
- `GET /api/signout` - Logout (requires authentication)

### Email Endpoints

- `GET /api/email/inbox` - Get inbox emails (requires authentication)
- `GET /api/email/:id` - Get a specific email (requires authentication)
- `POST /api/email/send` - Send an email (requires authentication)

For more detailed API documentation, see the [API Reference](docs/api/README.md).

## Documentation

Comprehensive documentation is available in the `docs` directory:

- [Getting Started](docs/getting-started/README.md): Installation and setup instructions
- [API Reference](docs/api/README.md): Detailed API documentation
- [Configuration](docs/config/README.md): Configuration options and examples
- [Architecture](docs/architecture/README.md): System architecture and component descriptions
- [Components](docs/components/README.md): Detailed component descriptions

## License

This project is licensed under the MIT License - see the LICENSE file for details.
