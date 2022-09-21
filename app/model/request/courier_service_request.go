package request

import (
	"encoding/json"
	"go-klikdokter/helper/global"
	"go-klikdokter/helper/message"
	"go-klikdokter/pkg/util/datatype"
	"time"

	validation "github.com/itgelo/ozzo-validation/v4"
)

// swagger:parameters SaveCourierServiceRequest
type ReqCourierServiceBody struct {
	//  in: body
	Body SaveCourierServiceRequest `json:"body"`
}

type SaveCourierServiceRequest struct {
	// Courier UId of the Courier Service
	// in: string
	CourierUId string `json:"courier_uid"`

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
	ETD_Min float64 `json:"etd_min"`

	// ETD Max of the Courier Service
	// in: float64
	ETD_Max float64 `json:"etd_max"`

	// Logo of the Courier Service
	// in: string
	Logo string `json:"logo"`

	// Cod Available of the Courier Service
	// in: integer
	CodAvailable int32 `json:"cod_available"`

	// Prescription Allowed of the Courier Service
	// in: integer
	PrescriptionAllowed int32 `json:"prescription_allowed"`

	// Cancelable of the Courier Service
	// in: integer
	Cancelable int32 `json:"cancelable"`

	// Tracking Available of the Courier Service
	// in: integer
	TrackingAvailable int32 `json:"tracking_available"`

	// Status of the Courier Service
	// in: integer
	Status *int32 `json:"status"`

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
	Insurance int32 `json:"insurance"`

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
	// example:"15:04:05+07"
	StartTime datatype.Time `json:"start_time"`

	// End Time of the Courier Service
	// example:"15:04:05+07"
	EndTime datatype.Time `json:"end_time"`

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

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`

	// Extend Jwt Info
	global.JWTInfo
}

// swagger:parameters courier service
type GetCourierServiceRequest struct {
	// name: id
	// in: path
	// required: true
	Id string `json:"id"`
}

// swagger:parameters CourierServiceListRequest
type CourierServiceListRequest struct {
	//Filter : {"courier_uid":["value","value"],"courier_type":["value","value"],"shipping_code":["value","value"],"shipping_name":["value","value"],"shipping_type_code":["value","value"],"status":[0,1]}
	// in: query
	Filter string `json:"filter"`

	// Maximun records per page
	// in: int32
	Limit int `schema:"limit" binding:"omitempty,numeric,min=1,max=100" json:"limit"`

	// Page No
	// in: int32
	Page int `schema:"page" binding:"omitempty,numeric,min=1" json:"page"`

	// Sort fields
	// in: string
	Sort string `schema:"sort" binding:"omitempty" json:"sort"`

	Filters CourierServiceListFilter `json:"-"`
}

type CourierServiceListFilter struct {
	CourierUID       []string `json:"courier_uid"`
	CourierType      []string `json:"courier_type"`
	ShippingCode     []string `json:"shipping_code"`
	ShippingName     []string `json:"shipping_name"`
	ShippingTypeCode []string `json:"shipping_type_code"`
	Status           []int    `json:"status"`
}

func (m *CourierServiceListRequest) GetFilter() {
	if len(m.Filter) > 0 {
		_ = json.Unmarshal([]byte(m.Filter), &m.Filters)
	}
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
	Uid string `json:"-" binding:"omitempty"`

	// CourierUId of the Courier Service
	// in: int
	CourierUId string `json:"courier_uid"`

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
	ETD_Min float64 `json:"etd_min"`

	// ETD Max of the Courier Service
	// in: float64
	ETD_Max float64 `json:"etd_max"`

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
	// example:"15:04:05+07"
	// in: time
	StartTime datatype.Time `json:"start_time"`

	// End Time of the Courier Service
	// example:"15:04:05+07"
	// in: time
	EndTime datatype.Time `json:"end_time"`

	// Repickup Fee of the Courier Service
	// in: float64
	Repickup float64 `json:"repickup"`

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`

	// Extend Jwt Info
	global.JWTInfo
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
		validation.Field(&req.CourierUId, validation.Required.Error(message.ErrReq.Message)),
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
