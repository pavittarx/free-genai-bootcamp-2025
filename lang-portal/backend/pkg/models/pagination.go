package models

// PaginationRequest represents the input parameters for pagination
type PaginationRequest struct {
	Page    int `json:"page" query:"page"`     // Current page number (1-indexed)
	Limit   int `json:"limit" query:"limit"`   // Number of items per page
	Offset  int `json:"offset" query:"offset"` // Calculated offset for database queries
	Sort    string `json:"sort" query:"sort"`  // Field to sort by
	Order   string `json:"order" query:"order"` // Sort order (asc/desc)
}

// PaginationResponse wraps paginated results with metadata
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// Validate checks and sets default values for pagination request
func (p *PaginationRequest) Validate() {
	// Set default page to 1 if not provided or invalid
	if p.Page < 1 {
		p.Page = 1
	}

	// Set default limit to 10 if not provided or out of range
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 10
	}

	// Calculate offset
	p.Offset = (p.Page - 1) * p.Limit

	// Set default sorting
	if p.Sort == "" {
		p.Sort = "id"
	}

	// Set default order
	if p.Order == "" {
		p.Order = "asc"
	}
}

// NewPaginationResponse creates a new pagination response
func NewPaginationResponse(data interface{}, total int64, req PaginationRequest) PaginationResponse {
	totalPages := int((total + int64(req.Limit) - 1) / int64(req.Limit))

	return PaginationResponse{
		Data:       data,
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
	}
}
