package pagination

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

// Params represents pagination parameters
type Params struct {
	Page     int
	PageSize int
	Offset   int
}

// Response represents pagination metadata for API responses
type Response struct {
	Page       int  `json:"page"`
	PageSize   int  `json:"page_size"`
	TotalItems int  `json:"total_items"`
	TotalPages int  `json:"total_pages"`
	HasMore    bool `json:"has_more"`
}

// DefaultPageSize is the default number of items per page
const DefaultPageSize = 20

// MaxPageSize is the maximum allowed page size
const MaxPageSize = 100

// GetParamsFromContext extracts pagination parameters from the request context
func GetParamsFromContext(c echo.Context) Params {
	// Get page parameter
	pageStr := c.QueryParam("page")
	page := 1
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Get page_size parameter
	pageSizeStr := c.QueryParam("page_size")
	pageSize := DefaultPageSize
	if pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
			if pageSize > MaxPageSize {
				pageSize = MaxPageSize
			}
		}
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	return Params{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// CreateResponse creates a pagination response based on the parameters and total items
func CreateResponse(params Params, totalItems int) Response {
	totalPages := (totalItems + params.PageSize - 1) / params.PageSize
	if totalPages < 1 {
		totalPages = 1
	}

	hasMore := params.Page < totalPages

	return Response{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}
}

// WrapResponse wraps the data with pagination metadata
func WrapResponse(data interface{}, pagination Response) map[string]interface{} {
	return map[string]interface{}{
		"data":       data,
		"pagination": pagination,
	}
}
