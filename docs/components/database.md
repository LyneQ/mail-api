# Database Component

The Database component is responsible for data persistence and retrieval in the MailAPI application.

## Overview

The Database component provides:
- Database connection management
- Data model definitions
- Schema migration
- Data access layer for other components

## Implementation

The Database component is implemented in the `db/database.go` file. It uses GORM (Go Object Relational Mapper) with SQLite as the database engine.

### Database Initialization

The database is initialized in the `Init` function:

```
func Init() {
    // Get database configuration from config.ini
    dbPath := config.GetDatabasePath()
    if dbPath == "" {
        dbPath = "./db/default.db" // Default path if not specified in config
    }

    db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
    if err != nil {
        _ = fmt.Errorf("failed to connect database: %v\n", err)
        log.Panicf("failed to connect database %v", err)
    }

    DB = db

    // Migrate the schema
    err = db.AutoMigrate(&User{})
    if err != nil {
        _ = fmt.Errorf("failed to migrate database: %v\n", err)
        return
    }

    fmt.Printf("[App] Database initialized at %s\n", dbPath)
}
```

This function:
1. Gets the database path from the configuration
2. Opens a connection to the SQLite database
3. Stores the database connection in a global variable `DB`
4. Automatically migrates the schema based on the defined models
5. Logs the successful initialization

### Data Models

Currently, the database has one main model:

#### User Model

```
type User struct {
    gorm.Model
    Username   string `json:"username" gorm:"uniqueIndex;not null"`
    Password   string `json:"-" gorm:"not null"`
    Role       string `json:"role" gorm:"not null;default:User"`
    IsVerified bool   `json:"is_verified" gorm:"default:false"`
}
```

The User model includes:
- `gorm.Model`: Embeds default fields (ID, CreatedAt, UpdatedAt, DeletedAt)
- `Username`: The user's unique identifier (typically an email address)
- `Password`: The user's password (not exposed in JSON)
- `Role`: The user's role (e.g., "User", "Admin")
- `IsVerified`: Whether the user's account has been verified

### Global Database Access

The database connection is exposed through a global variable:

```
var DB *gorm.DB
```

This allows other components to access the database without having to pass the connection around.

## Configuration

The Database component is configured through the `config/config.ini` file in the `[Database]` section:

```
[Database]
driver = sqlite3
path = ./db/dev.db
```

- **driver**: The database driver to use. Currently, only `sqlite3` is supported.
- **path**: The path to the SQLite database file. This can be relative to the application root or an absolute path.

## Usage Example

Here's an example of how the Database component is used in the application:

```
// Import the database package
import "github.com/lyneq/mailapi/db"

// Access the global DB variable
user := db.User{}
result := db.DB.Where("username = ?", "user@example.com").First(&user)
if result.Error != nil {
    // Handle error
}

// Create a new user
newUser := db.User{
    Username: "newuser@example.com",
    Password: hashedPassword,
    Role:     "User",
}
result = db.DB.Create(&newUser)
if result.Error != nil {
    // Handle error
}
```

## Schema Migration

The Database component uses GORM's AutoMigrate feature to automatically create and update the database schema based on the defined models. This happens during the `Init` function call.

The migration process:
1. Creates tables if they don't exist
2. Adds missing columns
3. Updates column types if needed
4. Creates indexes

This approach allows the schema to evolve as the models change, without requiring manual migration scripts.

## Security Considerations

The Database component includes several security features:

1. **Password Storage**: Passwords are not exposed in JSON responses
2. **Unique Usernames**: The `Username` field has a unique index to prevent duplicate accounts
3. **Soft Deletes**: The `gorm.Model` includes a `DeletedAt` field for soft deletes

## Limitations

The current implementation has a few limitations:

1. **Single Database Engine**: Only SQLite is supported, which may not be suitable for high-load production environments
2. **Limited Models**: Only the User model is defined, with no models for emails or other data
3. **No Connection Pooling**: The database connection is not pooled, which could limit performance under high load
4. **No Query Caching**: There's no caching of query results

## Future Enhancements

Potential improvements to the Database component include:

1. **Support for Multiple Database Engines**: Add support for PostgreSQL, MySQL, etc.
2. **Additional Models**: Add models for emails, attachments, etc.
3. **Connection Pooling**: Implement connection pooling for better performance
4. **Query Caching**: Add caching for frequently accessed data
5. **Database Migrations**: Implement proper database migrations for schema changes