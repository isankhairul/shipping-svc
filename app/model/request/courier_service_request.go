package request

import (
	"go-klikdokter/helper/message"
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters SaveCourierServiceRequest
type ReqCourierServiceBody struct {
	//  in: body
	Body SaveCourierServiceRequest `json:"body"`
}

type SaveCourierServiceRequest struct {
	// Uid of the Courier Service, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`

	// Courier Id of the Courier Service
	// in: int
	CourierId int `json:"courier_id"`

	// Courier Name of the Courier Service
	// in: string
	CourierName string `json:"courier_name"`

	// Shipping Code of the Courier Service
	// in: string
	ShippingCode string `json:"shipping_code"`

	// Shipping Name of the Courier Service
	// in: string
	ShippingName string `json:"shipping_name"`

	// Shipping Type of the Courier Service
	// in: string
	ShippingType string `json:"shipping_type"`

	// Shipping Description of the Courier Service
	// in: string
	ShippingDescription string `json:"shipping_description"`

	// ETD Min of the Courier Service
	// in: float64
	ETD_Min float64 `json:"ETD_Min"`

	// ETD Max of the Courier Service
	// in: float64
	ETD_Max float64 `json:"ETD_Max"`

	// Logo of the Courier Service
	// in: string
	Logo string `json:"logo"`

	// Cod Available of the Courier Service
	// in: integer
	CodAvailable int `json:"cod_available"`

	// Prescription Allowed of the Courier Service
	// in: integer
	PrescriptionAllowed int `json:"prescription_allowed"`

	// Cancelable of the Courier Service
	// in: integer
	Cancelable int `json:"cancelable"`

	// Tracking Available of the Courier Service
	// in: integer
	TrackingAvailable int `json:"tracking_available"`

	// Status of the Courier Service
	// in: integer
	Status int `json:"status"`

	// Max Weight of the Courier Service
	// in: float64
	MaxWeight float64 `json:"max_weight"`

	// Max Volume of the Courier Service
	// in: float64
	MaxVolume float64 `json:"max_volume"`

	// Max Distance of the Courier Service
	// in: float64
	MaxDistance float64 `json:"max_distance"`

	// Min Purchase of the Courier Service
	// in: integer
	MinPurchase float64 `json:"min_purchase"`

	// Max Purchase of the Courier Service
	// in: integer
	MaxPurchase float64 `json:"max_purchase"`

	// Insurance of the Courier Service
	// in: integer
	Insurance int `json:"insurance"`

	// Insurance Min of the Courier Service
	// in: float64
	InsuranceMin float64 `json:"insurance_min"`

	// Insurance Fee Type of the Courier Service
	// in: string
	InsuranceFeeType string `json:"insurance_fee_type"`

	// Insurance Fee of the Courier Service
	// in: float64
	InsuranceFee float64 `json:"insurance_fee"`

	// Start Time of the Courier Service
	// in: time
	StartTime time.Time `json:"start_time"`

	// End Time of the Courier Service
	// in: time
	EndTime time.Time `json:"end_time"`

	// Repickup Fee of the Courier Service
	// in: float64
	Repickup float64 `json:"repickup"`

	// Created At of the Courier Service
	// in: time
	CreatedAt time.Time `json:"created_at"`

	// Created By Type of the Courier Service
	// in: string
	CreatedBy string `json:"created_by"`

	// Updated At of the Courier Service
	// in: time
	UpdatedAt time.Time `json:"updated_at"`

	// Updated By Type of the Courier Service
	// in: string
	UpdatedBy string `json:"updated_by"`
}

// swagger:parameters courier service
type GetCourierServiceRequest struct {
	// name: id
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters list courier service
type CourierServiceListRequest struct {
	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1"`

	// Sort fields, example: name asc, uom desc
	// in: string
	Sort string `schema:"sort" binding:"omitempty"`

	// Shipping type
	// in: string
	ShippingType string `schema:"shipping_type" binding:"omitempty"`

	// Courier status
	// in: int
	Status int `schema:"status" binding:"omitempty"`
}

// swagger:parameters UpdateCourierServiceRequest
type ReqCourierServiceBodyUpdate struct {
	// Uid of the courier service.
	// in: path
	// required: true
	UId string `json:"uid"`
	//  in: body
	Body SaveCourierServiceRequest `json:"body"`
}

type UpdateCourierServiceRequest struct {
	// Uid of the courier service, use this on UPDATE function
	// in: string
	Uid string `json:"uid" binding:"omitempty"`

	// Courier Id of the Courier Service
	// in: int
	CourierId int `json:"courier_id"`

	// Courier Name of the Courier Service
	// in: string
	CourierName string `json:"courier_name"`

	// Shipping Code of the Courier Service
	// in: string
	ShippingCode string `json:"shipping_code"`

	// Shipping Name of the Courier Service
	// in: string
	ShippingName string `json:"shipping_name"`

	// Shipping Type of the Courier Service
	// in: string
	ShippingType string `json:"shipping_type"`

	// Shipping Description of the Courier Service
	// in: string
	ShippingDescription string `json:"shipping_description"`

	// ETD Min of the Courier Service
	// in: float64
	ETD_Min float64 `json:"ETD_Min"`

	// ETD Max of the Courier Service
	// in: float64
	ETD_Max float64 `json:"ETD_Max"`

	// Logo of the Courier Service
	// in: string
	Logo string `json:"logo"`

	// Cod Available of the Courier Service
	// in: integer
	CodAvailable int `json:"cod_available"`

	// Prescription Allowed of the Courier Service
	// in: integer
	PrescriptionAllowed int `json:"prescription_allowed"`

	// Cancelable of the Courier Service
	// in: integer
	Cancelable int `json:"cancelable"`

	// Tracking Available of the Courier Service
	// in: integer
	TrackingAvailable int `json:"tracking_available"`

	// Status of the Courier Service
	// in: integer
	Status int `json:"status"`

	// Max Weight of the Courier Service
	// in: float64
	MaxWeight float64 `json:"max_weight"`

	// Max Volume of the Courier Service
	// in: float64
	MaxVolume float64 `json:"max_volume"`

	// Max Distance of the Courier Service
	// in: float64
	MaxDistance float64 `json:"max_distance"`

	// Min Purchase of the Courier Service
	// in: integer
	MinPurchase float64 `json:"min_purchase"`

	// Max Purchase of the Courier Service
	// in: integer
	MaxPurchase float64 `json:"max_purchase"`

	// Insurance of the Courier Service
	// in: integer
	Insurance int `json:"insurance"`

	// Insurance Min of the Courier Service
	// in: float64
	InsuranceMin float64 `json:"insurance_min"`

	// Insurance Fee Type of the Courier Service
	// in: string
	InsuranceFeeType string `json:"insurance_fee_type"`

	// Insurance Fee of the Courier Service
	// in: float64
	InsuranceFee float64 `json:"insurance_fee"`

	// Start Time of the Courier Service
	// in: time
	StartTime time.Time `json:"start_time"`

	// End Time of the Courier Service
	// in: time
	EndTime time.Time `json:"end_time"`

	// Repickup Fee of the Courier Service
	// in: float64
	Repickup float64 `json:"repickup"`
}

// swagger:parameters CourierServiceRequestGetByUid
type CourierServiceGetByUid struct {
	// Uid of the Courier Service
	// in: path
	// required: true
	UId string `json:"uid"`
}

// swagger:parameters CourierServiceRequestDeleteByUid
type CourierServiceDeleteByUid struct {
	// Uid of the Courier Service
	// in: path
	// required: true
	UId string `json:"uid"`
}

func (req SaveCourierServiceRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.CourierId, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.CourierName, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.ShippingCode, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.ShippingName, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.ShippingType, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.ETD_Min, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.ETD_Max, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.MaxWeight, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.Logo, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.MaxVolume, validation.Required.Error(message.ErrReq.Message)),
		validation.Field(&req.MaxDistance, validation.Required.Error(message.ErrReq.Message)),
	)
}
