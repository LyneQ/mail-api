# API Reference

This document provides detailed information about the MailAPI endpoints, including request and response formats, authentication requirements, and examples.

## Authentication

Most endpoints in the MailAPI require authentication. The API uses session-based authentication.

### Authentication Flow

1. User signs up or signs in using the appropriate endpoints
2. The server sets a session cookie
3. Subsequent requests should include this cookie
4. Protected endpoints will check for a valid session

## API Endpoints

### Authentication Endpoints

#### Sign Up

Register a new user account.

- **URL**: `/api/signup`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "username": "user@example.com",
    "password": "securepassword"
  }
  ```
- **Success Response**: 
  - **Code**: 201 Created
  - **Content**: 
    ```json
    {
      "user": {
        "id": 1,
        "username": "user@example.com",
        "role": "User",
        "is_verified": false
      }
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request
  - **Content**: 
    ```json
    {
      "error": "Username already exists"
    }
    ```

#### Sign In

Authenticate a user and create a session.

- **URL**: `/api/signin`
- **Method**: `POST`
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "username": "user@example.com",
    "password": "securepassword"
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "user": {
        "id": 1,
        "username": "user@example.com",
        "role": "User",
        "is_verified": false
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**: 
    ```json
    {
      "error": "Invalid credentials"
    }
    ```

#### Get Current User

Get information about the currently authenticated user.

- **URL**: `/api/me`
- **Method**: `GET`
- **Auth Required**: Yes
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "user": {
        "id": 1,
        "username": "user@example.com",
        "role": "User",
        "is_verified": false
      }
    }
    ```
- **Error Response**:
  - **Code**: 401 Unauthorized
  - **Content**: 
    ```json
    {
      "error": "Not authenticated"
    }
    ```

#### Sign Out

End the current user session.

- **URL**: `/api/signout`
- **Method**: `GET`
- **Auth Required**: Yes
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Successfully signed out"
    }
    ```

### Email Endpoints

#### Get Inbox

Retrieve emails from the user's inbox.

- **URL**: `/api/email/inbox`
- **Method**: `GET`
- **Auth Required**: Yes
- **Query Parameters**:
  - `limit` (optional): Maximum number of emails to retrieve (default: 50)
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "emails": [
        {
          "id": "1",
          "from": "sender@example.com",
          "to": ["recipient@example.com"],
          "subject": "Hello",
          "date": "2023-01-01T12:00:00Z"
        },
        {
          "id": "2",
          "from": "another@example.com",
          "to": ["recipient@example.com"],
          "subject": "Meeting",
          "date": "2023-01-02T14:30:00Z"
        }
      ]
    }
    ```

#### Get Email by ID

Retrieve a specific email with full details.

- **URL**: `/api/email/:id`
- **Method**: `GET`
- **Auth Required**: Yes
- **URL Parameters**:
  - `id`: ID of the email to retrieve
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "email": {
        "id": "1",
        "from": "sender@example.com",
        "to": ["recipient@example.com"],
        "subject": "Hello",
        "body": "<p>This is the email body</p>",
        "date": "2023-01-01T12:00:00Z",
        "attachments": [
          {
            "filename": "document.pdf",
            "mime_type": "application/pdf"
          }
        ]
      }
    }
    ```
- **Error Response**:
  - **Code**: 404 Not Found
  - **Content**: 
    ```json
    {
      "error": "Email not found"
    }
    ```

#### Send Email

Send a new email.

- **URL**: `/api/email/send`
- **Method**: `POST`
- **Auth Required**: Yes
- **Request Body**:
  ```json
  {
    "to": ["recipient@example.com"],
    "subject": "Hello",
    "body": "<p>This is the email body</p>",
    "attachments": [
      {
        "filename": "document.pdf",
        "content": "base64_encoded_content",
        "mime_type": "application/pdf"
      }
    ]
  }
  ```
- **Success Response**: 
  - **Code**: 200 OK
  - **Content**: 
    ```json
    {
      "message": "Email sent successfully"
    }
    ```
- **Error Response**:
  - **Code**: 400 Bad Request
  - **Content**: 
    ```json
    {
      "error": "Invalid recipient email address"
    }
    ```

## Error Handling

The API uses standard HTTP status codes to indicate the success or failure of a request:

- `200 OK`: The request was successful
- `201 Created`: A resource was successfully created
- `400 Bad Request`: The request was malformed or invalid
- `401 Unauthorized`: Authentication is required or failed
- `404 Not Found`: The requested resource was not found
- `500 Internal Server Error`: An unexpected error occurred on the server

Error responses include a JSON object with an `error` field containing a description of the error.