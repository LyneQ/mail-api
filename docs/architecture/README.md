# Architecture Overview

This document provides an overview of the MailAPI architecture, including its components, their interactions, and the overall system design.

## System Architecture

MailAPI follows a layered architecture pattern with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                        API Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐  │
│  │  Auth Controller│  │ Email Controller│  │    ...      │  │
│  └─────────────────┘  └─────────────────┘  └─────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                     Session Layer                           │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │   SMTP Client   │  │   IMAP Client   │                   │
│  └─────────────────┘  └─────────────────┘                   │
├─────────────────────────────────────────────────────────────┤
│                    Data Layer                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                     Database                        │    │
│  └─────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### API Layer

The API layer is responsible for handling HTTP requests and responses. It uses the Echo web framework to define routes and controllers.

Key components:
- **Auth Controller**: Handles user authentication (signup, signin, signout)
- **Email Controller**: Handles email operations (inbox, get email, send email)

### Session Layer

The session layer manages user sessions and authentication state. It provides middleware for protecting routes that require authentication.

Key components:
- **Session Management**: Handles session creation, validation, and destruction
- **Authentication Middleware**: Verifies user authentication for protected routes

### Service Layer

The service layer contains the business logic for interacting with external services like email servers.

Key components:
- **SMTP Client**: Handles sending emails via SMTP
- **IMAP Client**: Handles retrieving emails via IMAP

### Data Layer

The data layer is responsible for data persistence and retrieval.

Key components:
- **Database**: SQLite database for storing user information
- **Models**: GORM models for database entities

## Component Interactions

### Authentication Flow

1. User sends credentials to `/api/signin`
2. Auth controller validates credentials against the database
3. If valid, session layer creates a new session
4. Session cookie is returned to the client
5. Subsequent requests include the session cookie
6. Authentication middleware validates the session for protected routes

### Email Sending Flow

1. Authenticated user sends email data to `/api/email/send`
2. Email controller validates the request
3. SMTP client connects to the configured SMTP server
4. Email is sent via the SMTP server
5. Response is returned to the client

### Email Retrieval Flow

1. Authenticated user requests inbox at `/api/email/inbox`
2. Email controller processes the request
3. IMAP client connects to the configured IMAP server
4. Emails are retrieved from the server
5. Email data is returned to the client

## Key Design Decisions

### Web Framework

The application uses the Echo web framework for its API layer. Echo was chosen for its:
- Lightweight design
- High performance
- Middleware support
- Simple and intuitive API

### Database

SQLite with GORM is used for data persistence. This combination was chosen for:
- Simplicity of setup (no separate database server required)
- GORM's object-relational mapping capabilities
- Ease of development and testing

### Email Handling

The application uses separate libraries for SMTP and IMAP:
- **gomail.v2** for SMTP operations (sending emails)
- **go-imap** for IMAP operations (retrieving emails)

This separation allows for specialized handling of each protocol while maintaining a clean architecture.

### Authentication

The application uses session-based authentication rather than token-based authentication (like JWT). This decision was made for:
- Simplicity of implementation
- Built-in session management in Echo
- Ease of session revocation

## Security Considerations

The architecture includes several security measures:

1. **Password Storage**: User passwords are not stored in plain text (they should be hashed before storage)
2. **Session Management**: Sessions have expiration times and can be revoked
3. **CORS Protection**: Only allowed domains can make cross-origin requests
4. **TLS for Email**: Connections to email servers use TLS encryption

## Scalability

The current architecture can be scaled in several ways:

1. **Horizontal Scaling**: Deploy multiple instances behind a load balancer
2. **Database Scaling**: Replace SQLite with a more scalable database like PostgreSQL
3. **Caching**: Add a caching layer for frequently accessed data
4. **Asynchronous Processing**: Move email sending to a background job queue

## Future Enhancements

Potential architectural improvements include:

1. **Microservices**: Split the application into separate services for auth and email
2. **API Gateway**: Add an API gateway for rate limiting and additional security
3. **Event-Driven Architecture**: Use message queues for asynchronous processing
4. **Containerization**: Package the application as Docker containers for easier deployment