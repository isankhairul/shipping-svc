package response

import "go-klikdokter/app/model/entity"

// swagger:model ChannelCourierPaginationReponse
type ChannelCourierPaginationReponse struct {
	Records      int64 `json:"records"`
	TotalRecords int64 `json:"total_records"`
	Limit        int   `json:"limit"`
	Page         int   `json:"page"`
	TotalPage    int   `json:"total_page"`
	Items        []*entity.ChannelCourierDTO
}

// swagger:response ChannelCourierDTO
type ChannelCourierDTOResponse struct {
	// in: body
	Response entity.ChannelCourierDTO
}
