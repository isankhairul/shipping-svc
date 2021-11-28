package response

type PaginationResponse struct {
	Records      int64 `json:"records"`
	TotalRecords int64 `json:"total_records"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	TotalPage    int   `json:"total_page"`
}

func (p *PaginationResponse) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *PaginationResponse) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *PaginationResponse) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}
