# Authentication Component

The Authentication component is responsible for user registration, login, session management, and access control in the MailAPI application.

## Overview

The Authentication component provides a secure way to:
- Register new users
- Authenticate users
- Maintain user sessions
- Protect routes that require authentication

## Implementation

The Authentication component is implemented across several files:
- `api/auth/controller.go`: Defines the authentication routes
- `api/auth/view.go`: Implements the handlers for authentication requests
- `internal/session/main.go`: Manages user sessions
- `internal/middleware/middleware.go`: Provides middleware for protecting routes

### User Model

The User model is defined in `db/database.go`:

```
type User struct {
    gorm.Model
    Username   string `json:"username" gorm:"uniqueIndex;not null"`
    Password   string `json:"-" gorm:"not null"`
    Role       string `json:"role" gorm:"not null;default:User"`
    IsVerified bool   `json:"is_verified" gorm:"default:false"`
}
```

Key fields:
- `Username`: The user's unique identifier (typically an email address)
- `Password`: The user's password (stored securely and not exposed in JSON)
- `Role`: The user's role (e.g., "User", "Admin")
- `IsVerified`: Whether the user's account has been verified

### Authentication Routes

The authentication routes are defined in `api/auth/controller.go`:

1. **Sign Up** (`POST /api/signup`): Register a new user
2. **Sign In** (`POST /api/signin`): Authenticate a user and create a session
3. **Me** (`GET /api/me`): Get information about the currently authenticated user
4. **Sign Out** (`GET /api/signout`): End the current user session

### Session Management

Sessions are managed by the code in `internal/session/main.go`. The session system:
1. Creates a new session when a user signs in
2. Stores session data in the database
3. Sets a session cookie in the user's browser
4. Validates the session on subsequent requests
5. Destroys the session when the user signs out

### Authentication Middleware

The `RequireAuth` middleware in `internal/middleware/middleware.go` protects routes that require authentication:

1. It checks for a valid session
2. If a valid session exists, it allows the request to proceed
3. If no valid session exists, it returns a 401 Unauthorized response

## Authentication Flow

### Registration Flow

1. Client sends a POST request to `/api/signup` with username and password
2. Server validates the input
3. Server checks if the username already exists
4. If validation passes, the server creates a new user with a hashed password
5. Server returns the created user (without the password)

### Login Flow

1. Client sends a POST request to `/api/signin` with username and password
2. Server validates the input
3. Server checks if the username exists and the password is correct
4. If authentication is successful, the server creates a new session
5. Server sets a session cookie and returns the user information

### Authentication Check Flow

1. Client sends a GET request to `/api/me` with the session cookie
2. Server validates the session
3. If the session is valid, the server returns the user information
4. If the session is invalid, the server returns a 401 Unauthorized response

### Logout Flow

1. Client sends a GET request to `/api/signout` with the session cookie
2. Server destroys the session
3. Server clears the session cookie
4. Server returns a success message

## Security Considerations

The Authentication component includes several security features:

1. **Password Hashing**: Passwords should be hashed before storage (not stored in plain text)
2. **Session Expiration**: Sessions have a limited lifetime
3. **HTTPS Only Cookies**: Session cookies should be sent only over HTTPS
4. **CSRF Protection**: The API includes CSRF protection for session cookies

## Configuration

The Authentication component doesn't have specific configuration options in the `config.ini` file. However, session behavior can be configured when initializing the session system:

```
session.Init(db.DB, false)
```

The second parameter (`false`) indicates whether to use secure cookies (HTTPS only). In production, this should be set to `true`.

## Limitations

The current implementation has a few limitations:

1. **No Password Reset**: There's no built-in mechanism for password reset
2. **No Multi-factor Authentication**: The system only supports single-factor authentication
3. **No Account Verification**: While there's an `IsVerified` field, there's no implementation for verifying accounts
4. **No Role-based Access Control**: While there's a `Role` field, there's no implementation for role-based permissions

These limitations could be addressed in future versions of the component.