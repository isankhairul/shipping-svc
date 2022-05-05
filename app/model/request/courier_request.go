package request

// swagger:parameters SaveCourierRequest
type ReqCourierBody struct {
	//  in: body
	Body SaveCourierRequest `json:"body"`
}

type SaveCourierRequest struct {
	// Name of the courier
	// in: string
	CourierName string `json:"name"`

	// Code of the courier
	// in: string
	Code string `json:"code"`

	// Uid of the courỉe, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`

	// type of courier
	// in: string
	CourierType string `json:"courier_type"`
}

// swagger:parameters courier
type GetCourierRequest struct {
	// name: id
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters list courier
type CourierListRequest struct {
	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: string
	Sort string `schema:"sort" binding:"omitempty"`

	// Name of the Courier
	// in: string
	CourierName string `json:"courier_name"`

	// Gender of the Courier
	// in: string
	Code string `json:"code"`

	// Name of Courier Type
	// in: string
	// required: false
	CourierType string `json:"courier_type"`

	// Name of Courier Type
	// in: string
	// required: false
	Logo string `json:"logo"`

	// Name of Courier Type
	// in: string
	// required: false
	HidePurpose int `json:"hide_purpose"`

	// Name of Courier Type
	// in: string
	// required: false
	CourierApiIntegration int `json:"courier_api_intergration"`

	// Name of Courier Type
	// in: string
	// required: false
	UseGeocoodinate int `json:"use_geocoodinate"`

	// Name of Courier Type
	// in: string
	// required: false
	ProvideAirwaybill int `json:"provide_airwaybill"`

	// Status of Courier
	// in: integer
	// required: false
	Status int `json:"status"`
}
