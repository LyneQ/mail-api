# Configuration Guide

This document provides detailed information about configuring the MailAPI application.

## Configuration File

MailAPI uses an INI-style configuration file located at `config/config.ini`. This file contains all the necessary settings for the application to function properly.

## Configuration Sections

The configuration file is divided into several sections, each controlling a different aspect of the application.

### AllowedDomains

This section defines the domains that are allowed to make cross-origin requests to the API.

```ini
[AllowedDomains]
domains = localhost, example.com, yourdomain.com
```

- **domains**: A comma-separated list of domain names that are allowed to access the API. These domains will be used to configure CORS (Cross-Origin Resource Sharing) headers.

### Database

This section configures the database connection.

```ini
[Database]
driver = sqlite3
path = ./db/dev.db
```

- **driver**: The database driver to use. Currently, only `sqlite3` is supported.
- **path**: The path to the SQLite database file. This can be relative to the application root or an absolute path.

### Api

This section configures the API server.

```ini
[Api]
port = 1323
```

- **port**: The port number on which the API server will listen for incoming connections.

### SMTP

This section configures the SMTP client for sending emails.

```ini
[SMTP]
host = smtp.example.com
port = 587
username = your_username
password = your_password
```

- **host**: The hostname of the SMTP server.
- **port**: The port number of the SMTP server (typically 587 for TLS or 465 for SSL).
- **username**: The username for authenticating with the SMTP server.
- **password**: The password for authenticating with the SMTP server.

### IMAP

This section configures the IMAP client for receiving emails.

```ini
[IMAP]
host = imap.example.com
port = 993
username = your_username
password = your_password
```

- **host**: The hostname of the IMAP server.
- **port**: The port number of the IMAP server (typically 993 for SSL).
- **username**: The username for authenticating with the IMAP server.
- **password**: The password for authenticating with the IMAP server.

## Environment-Specific Configuration

You can create different configuration files for different environments:

- `config/config.development.ini` for development
- `config/config.production.ini` for production
- `config/config.test.ini` for testing

To use a specific configuration file, set the `APP_ENV` environment variable:

```bash
export APP_ENV=production
```

If `APP_ENV` is not set, the application will use `config/config.ini` by default.

## Sensitive Information

Be careful with sensitive information like database credentials and email passwords. For production environments, consider:

1. Using environment variables instead of hardcoding values in the config file
2. Setting appropriate file permissions on the config file
3. Using a secrets management solution

## Configuration Loading

The configuration is loaded when the application starts. If there are any errors in the configuration file, the application will log an error and exit.

You can see the configuration loading code in `config/config.go`.

## Example Configuration

Here's a complete example of a configuration file:

```ini
[AllowedDomains]
domains = localhost, example.com

[Database]
driver = sqlite3
path = ./db/dev.db

[Api]
port = 1323

[SMTP]
host = smtp.gmail.com
port = 587
username = your_email@gmail.com
password = your_app_password

[IMAP]
host = imap.gmail.com
port = 993
username = your_email@gmail.com
password = your_app_password
```

This configuration sets up the application to:
- Allow requests from localhost and example.com
- Use a SQLite database at ./db/dev.db
- Run the API server on port 1323
- Connect to Gmail's SMTP and IMAP servers for email functionality