package response

// swagger:model ShippmentPredefinedDetail
type ShippmentPredefined struct {
	UID string `json:"uid"`

	Type string `json:"type"`

	Code string `json:"code"`

	Title string `json:"title"`

	Note string `json:"note"`

	Status int32 `json:"status"`
}

// swagger:response GetShippmentPredefined
type GetShippmentPredefinedResponse struct {
	//in: body
	Response ShippmentPredefined `json:"response"`
}
