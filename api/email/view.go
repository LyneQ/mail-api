package email

import (
	"fmt"
	"github.com/lyneq/mailapi/config"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/internal/smtpClient"
)

// EmailResponse represents the response structure for email data
type EmailResponse struct {
	ID          string       `json:"id"`
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	Date        string       `json:"date"`
	Body        string       `json:"body,omitempty"`
	Labels      []string     `json:"labels"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents an email attachment in the response
type Attachment struct {
	Filename string `json:"filename"`
	MimeType string `json:"mime_type"`
	Size     int    `json:"size"`
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

	imapClient := smtpclient.NewIMAPClientFromConfig()

	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	messages, err := imapClient.GetInbox(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get inbox: %v", err),
		})
	}

	var emails []EmailResponse
	for _, msg := range messages {
		email := EmailResponse{
			ID:      msg.ID,
			From:    msg.From,
			To:      msg.To,
			Subject: msg.Subject,
			Date:    msg.Date.Format("2006-01-02 15:04:05"),
			Labels:  msg.Flags,
		}

		emails = append(emails, email)
	}

	folders, err := imapClient.GetFolders()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get folders: %v", err),
		})
	}

	return c.JSON(http.StatusOK, Response{Folders: folders, Emails: emails})
}

func getFolderView(c echo.Context) error {
	imapClient := smtpclient.NewIMAPClientFromConfig()

	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	limitStr := c.QueryParam("limit")
	limit := 20
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

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

	var response []EmailResponse
	for _, msg := range messages {
		email := EmailResponse{
			ID:      msg.ID,
			From:    msg.From,
			To:      msg.To,
			Subject: msg.Subject,
			Date:    msg.Date.Format("2006-01-02 15:04:05"),
			Labels:  msg.Flags,
		}

		response = append(response, email)
	}

	return c.JSON(http.StatusOK, response)
}

// getEmailView handles the request to get a specific email by ID
func getEmailView(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email ID is required",
		})
	}

	imapClient := smtpclient.NewIMAPClientFromConfig()

	// Connect to IMAP server
	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	folder := c.QueryParam("folder")

	message, err := imapClient.GetEmailByID(id, folder)
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
		Labels:  message.Flags,
	}

	for _, att := range message.Attachments {
		response.Attachments = append(response.Attachments, Attachment{
			Filename: att.Filename,
			MimeType: att.MimeType,
			Size:     len(att.Content),
		})
	}

	return c.JSON(http.StatusOK, response)
}

// getFoldersView handles the request to get all mail folders
func getFoldersView(c echo.Context) error {
	imapClient := smtpclient.NewIMAPClientFromConfig()

	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	folders, err := imapClient.GetFolders()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get folders: %v", err),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"folders": folders,
	})
}

// sendEmailView handles the request to send an email
func sendEmailView(c echo.Context) error {
	req := new(SendEmailRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Invalid request: %v", err),
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Validation error: %v", err),
		})
	}

	smtpClient := smtpclient.NewSMTPClientFromConfig()

	if err := smtpClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to SMTP server: %v", err),
		})
	}

	sender := config.GetIMAPConfig().Username

	err := smtpClient.SendMessage(
		sender,
		req.To,
		req.Subject,
		req.Body,
		nil,
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
