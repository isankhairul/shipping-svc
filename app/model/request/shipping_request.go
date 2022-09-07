package request

//swagger:parameters ShippingRate
type GetShippingRate struct {
	//in: body
	Body GetShippingRateRequest `json:"body"`
}

//swagger:parameters ShippingRateByShippingType
type GetShippingRateByShippingType struct {
	//in: path
	ShippingType string `json:"shipping-type"`
	//in: body
	Body GetShippingRateRequest `json:"body"`
}
type GetShippingRateRequest struct {
	ShippingType        string
	ChannelUID          string            `json:"channel_uid"`
	TotalWeight         float64           `json:"total_weight"`
	TotalWidth          float64           `json:"total_width"`
	TotalHeight         float64           `json:"total_heigth"`
	TotalLength         float64           `json:"total_length"`
	TotalProductPrice   float64           `json:"total_product_price"`
	ContainPrescription bool              `json:"contain_prescription"`
	Origin              AreaDetailPayload `json:"origin"`
	Destination         AreaDetailPayload `json:"destination"`
	CourierServiceUID   []string          `json:"courier_service_uid"`
}

type AreaDetailPayload struct {
	CountryCode string `json:"country_code"`
	PostalCode  string `json:"postal_code"`
	Subdistrict string `json:"subdistrict"`
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
}

type ChannelCourierServicePayloadItem struct {
	CourierServiceUID string `json:"courier_service_uid"`
}
