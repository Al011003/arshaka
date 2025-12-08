package response

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalRows  int `json:"total_rows"`
	TotalPages int `json:"total_pages"`
}

type PaginatedResponse struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}