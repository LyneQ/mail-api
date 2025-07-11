package utils

import (
	"encoding/base64"
	"fmt"
	"strings"

	smtpclient "github.com/lyneq/mailapi/internal/smtpClient"
)

func AttachmentToHTML(attachment smtpclient.Attachment, size int) string {
	switch {
	case strings.HasPrefix(attachment.MimeType, "image/"):
		encoded := base64.StdEncoding.EncodeToString(attachment.Content)
		return fmt.Sprintf(`<img src="data:%s;base64,%s" alt="%s" style="max-width: 100%%; display: block;" />`,
			attachment.MimeType, encoded, attachment.Filename)

	case strings.HasPrefix(attachment.MimeType, "text/"):
		limited := attachment.Content
		if size > 0 && len(limited) > size {
			limited = limited[:size]
		}
		return fmt.Sprintf(`<pre>%s</pre>`, string(limited))

	case strings.HasPrefix(attachment.MimeType, "application/pdf"):
		encoded := base64.StdEncoding.EncodeToString(attachment.Content)
		return fmt.Sprintf(`<embed src="data:%s;base64,%s" type="application/pdf" width="100%%" height="600px" />`,
			attachment.MimeType, encoded)

	default:
		encoded := base64.StdEncoding.EncodeToString(attachment.Content)
		return fmt.Sprintf(`<a href="data:%s;base64,%s" download="%s">Download %s</a>`,
			attachment.MimeType, encoded, attachment.Filename, attachment.Filename)
	}
}
