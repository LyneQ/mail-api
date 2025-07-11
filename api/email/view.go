package email

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lyneq/mailapi/config"
	"github.com/lyneq/mailapi/internal/pagination"
	"github.com/lyneq/mailapi/internal/smtpClient"
	"github.com/lyneq/mailapi/internal/utils"
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
	HTML     string `json:"html,omitempty"`
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

// getInboxView handles the request to get the user's inbox with pagination
func getInboxView(c echo.Context) error {
	imapClient := smtpclient.NewIMAPClientFromConfig()

	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	paginationParams := pagination.GetParamsFromContext(c)

	limitStr := c.QueryParam("limit")
	if limitStr != "" && paginationParams.PageSize == pagination.DefaultPageSize {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			paginationParams.PageSize = parsedLimit
			if paginationParams.PageSize > pagination.MaxPageSize {
				paginationParams.PageSize = pagination.MaxPageSize
			}
		}
	}

	result, err := imapClient.GetInbox(paginationParams.Page, paginationParams.PageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get inbox: %v", err),
		})
	}

	var emails []EmailResponse
	for _, msg := range result.Messages {
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

	paginationResponse := pagination.CreateResponse(paginationParams, int(result.TotalCount))

	return c.JSON(http.StatusOK, map[string]interface{}{
		"folders":    folders,
		"emails":     emails,
		"pagination": paginationResponse,
	})
}

// getFolderView handles the request to get messages from a specific folder with pagination
func getFolderView(c echo.Context) error {
	imapClient := smtpclient.NewIMAPClientFromConfig()

	if err := imapClient.Connect(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to connect to IMAP server: %v", err),
		})
	}
	defer imapClient.Disconnect()

	paginationParams := pagination.GetParamsFromContext(c)

	limitStr := c.QueryParam("limit")
	if limitStr != "" && paginationParams.PageSize == pagination.DefaultPageSize {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			paginationParams.PageSize = parsedLimit
			if paginationParams.PageSize > pagination.MaxPageSize {
				paginationParams.PageSize = pagination.MaxPageSize
			}
		}
	}

	folderName := c.QueryParam("name")
	if folderName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Folder name is required",
		})
	}

	result, err := imapClient.GetFolderMessages(folderName, paginationParams.Page, paginationParams.PageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to get folder messages: %v", err),
		})
	}

	var emails []EmailResponse
	for _, msg := range result.Messages {
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

	paginationResponse := pagination.CreateResponse(paginationParams, int(result.TotalCount))

	return c.JSON(http.StatusOK, pagination.WrapResponse(emails, paginationResponse))
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

	response := EmailResponse{
		ID:      message.ID,
		From:    message.From,
		To:      message.To,
		Subject: message.Subject,
		Date:    message.Date.Format("2006-01-02 15:04:05"),
		Labels:  message.Flags,
	}

	charLimit := 10000

	if message.Size > 0 {

		if message.Size > 1000000 {
			charLimit = 5000
		} else if message.Size > 500000 {
			charLimit = 7500
		} else {

			charLimit = 15000
		}
	}

	bodyToClean := message.Body
	if len(bodyToClean) > charLimit {
		bodyToClean = bodyToClean[:charLimit]
	}
	emailBody := CleanBinaryData(bodyToClean)

	for _, att := range message.Attachments {
		attachmentLimit := 1000

		// Use both the message size and attachment size to determine the character limit
		// This helps optimize the cleaning process for large messages with large attachments
		if message.Size > 1000000 { // If message is larger than 1MB
			// For very large messages, be more aggressive with attachment limits
			attachmentLimit = 500
		} else if len(att.Content) > 100000 {
			attachmentLimit = 750
		} else if message.Size > 500000 {
			attachmentLimit = 800
		}

		attachment := Attachment{
			Filename: att.Filename,
			MimeType: att.MimeType,
			Size:     attachmentLimit,
			HTML:     utils.AttachmentToHTML(att, attachmentLimit),
		}

		response.Attachments = append(response.Attachments, attachment)

		emailBody += "\n\n" + attachment.HTML
	}

	response.Body = emailBody

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

// CleanBinaryData removes binary data and non-printable characters from the email body
// This function uses a whitelist approach to keep only valid text characters
func CleanBinaryData(body string) string {
	reBase64Tags := regexp.MustCompile(`(?s)<(img|embed|object)[^>]*base64[^>]*>`)
	body = reBase64Tags.ReplaceAllString(body, "")

	rePureBase64 := regexp.MustCompile(`data:[^;]+;base64,[a-zA-Z0-9+/=]+`)
	body = rePureBase64.ReplaceAllString(body, "")

	reBinaryGarbage := regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F-\x9F\xAD�]+`)

	body = reBinaryGarbage.ReplaceAllString(body, "")

	var b strings.Builder
	for _, r := range body {
		switch {
		case r >= 32 && r <= 126: // ASCII printable
			b.WriteRune(r)
		case r == '\n' || r == '\r' || r == '\t':
			b.WriteRune(r)
		case r >= 0x00A0 && r <= 0x00FF: // Latin-1 Supplement
			b.WriteRune(r)
		case r >= 0x0100 && r <= 0x017F: // Latin Extended-A
			b.WriteRune(r)
		default:
			// On skip tout ce qui est trop chelou
			continue
		}
	}

	clean := b.String()

	reCollapse := regexp.MustCompile(`[\s]{2,}`)
	clean = reCollapse.ReplaceAllString(clean, " ")

	return strings.TrimSpace(clean)
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
