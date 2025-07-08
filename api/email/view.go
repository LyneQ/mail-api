package email

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/internal/smtpClient"
)

// EmailResponse represents the response structure for email data
type EmailResponse struct {
	ID      string   `json:"id"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Date    string   `json:"date"`
	Body    string   `json:"body,omitempty"`
	Flags   []string `json:"flags"`
}

type FolderResponse struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Response struct {
	Folders any             `json:"folders"`
	Emails  []EmailResponse `json:"inbox"`
}

// SendEmailRequest represents the request structure for sending an email
type SendEmailRequest struct {
	To       []string `json:"to" validate:"required,min=1,dive,email"`
	Subject  string   `json:"subject" validate:"required"`
	Body     string   `json:"body" validate:"required"`
	HTMLBody bool     `json:"html_body"`
}

// getInboxView handles the request to get the user's inbox
func getInboxView(c echo.Context) error {
	// Create IMAP client
	imapClient := smtpclient.NewIMAPClientFromConfig()

	// Connect to IMAP server
	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	// Get limit parameter, default to 20 if not provided
	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get inbox messages
	messages, err := imapClient.GetInbox(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get inbox: %v", err),
		})
	}

	// Convert to response format
	var emails []EmailResponse
	for _, msg := range messages {
		emails = append(emails, EmailResponse{
			ID:      msg.ID,
			From:    msg.From,
			To:      msg.To,
			Subject: msg.Subject,
			Date:    msg.Date.Format("2006-01-02 15:04:05"),
			Flags:   msg.Flags,
		})
	}

	folders, err := imapClient.GetFolders()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get folders: %v", err),
		})
	}

	// Add folders info if needed
	// For now, we're just returning the basic email details including the body'

	return c.JSON(http.StatusOK, Response{Folders: folders, Emails: emails})
}

func getFolderView(c echo.Context) error {
	// Create IMAP client
	imapClient := smtpclient.NewIMAPClientFromConfig()

	// Connect to IMAP server
	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	// Get limit parameter, default to 20 if not provided
	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Get inbox messages
	folderName := c.QueryParam("name")
	if folderName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Folder name is required",
		})
	}
	fmt.Println(folderName)
	messages, err := imapClient.GetFolderMessages(folderName, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get inbox: %v", err),
		})
	}

	// Convert to response format
	var response []EmailResponse
	for _, msg := range messages {
		response = append(response, EmailResponse{
			ID:      msg.ID,
			From:    msg.From,
			To:      msg.To,
			Subject: msg.Subject,
			Date:    msg.Date.Format("2006-01-02 15:04:05"),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// getEmailView handles the request to get a specific email by ID
func getEmailView(c echo.Context) error {
	// Get email ID from path parameter
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email ID is required",
		})
	}

	// Create IMAP client
	imapClient := smtpclient.NewIMAPClientFromConfig()

	// Connect to IMAP server
	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	// Get the specific email by ID
	message, err := imapClient.GetEmailByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get email: %v", err),
		})
	}

	// Convert to response format
	response := EmailResponse{
		ID:      message.ID,
		From:    message.From,
		To:      message.To,
		Subject: message.Subject,
		Date:    message.Date.Format("2006-01-02 15:04:05"),
		Body:    message.Body,
	}

	// Add attachments info if needed
	// For now, we're just returning the basic email details including the body

	return c.JSON(http.StatusOK, response)
}

// sendEmailView handles the request to send an email
func sendEmailView(c echo.Context) error {
	// Parse request body
	req := new(SendEmailRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Validation error: %v", err),
		})
	}

	// Create SMTP client
	smtpClient := smtpclient.NewSMTPClientFromConfig()

	// Connect to SMTP server
	if err := smtpClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to SMTP server: %v", err),
		})
	}

	// Get sender from configuration
	sender := "Contact@lyneq.tech" // Using the email from config.ini

	// Send email
	err := smtpClient.SendMessage(
		sender,
		req.To,
		req.Subject,
		req.Body,
		nil, // No attachments for now
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to send email: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Email sent successfully",
	})
}
