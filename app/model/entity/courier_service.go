package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/pkg/util/datatype"
)

type CourierServiceDetailDTO struct {
	Uid string `gorm:"not null" json:"uid"`

	// Courier UId of the Courier Service
	// in: string
	CourierName string `gorm:"not null" json:"courier_name"`

	// Courier UId of the Courier Service
	// in: string
	CourierType string `gorm:"not null" json:"courier_type"`

	// Courier UId of the Courier Service
	// in: string
	CourierUId string `gorm:"not null" json:"courier_uid"`

	// Shipping Code of the Courier Service
	// in: string
	ShippingCode string `gorm:"not null" json:"shipping_code"`

	// Shipping Name of the Courier Service
	// in: string
	ShippingName string `gorm:"not null" json:"shipping_name"`

	// Shipping Type of the Courier Service
	// in: string
	ShippingType string `gorm:"not null" json:"shipping_type"`

	// Shipping Description of the Courier Service
	// in: string
	ShippingDescription string `json:"shipping_description"`

	// ETD Min of the Courier Service
	// in: float64
	ETD_Min float64 `gorm:"not null" json:"etd_min"`

	// ETD Max of the Courier Service
	// in: float64
	ETD_Max float64 `gorm:"not null" json:"etd_max"`

	// Logo of the Courier Service
	// in: string
	Logo string `gorm:"not null" json:"logo"`

	// Cod Available of the Courier Service
	// in: integer
	CodAvailable int `gorm:"not null;default:0" json:"cod_available"`

	// Prescription Allowed of the Courier Service
	// in: integer
	PrescriptionAllowed int `gorm:"not null;default:0" json:"prescription_allowed"`

	// Cancelable of the Courier Service
	// in: integer
	Cancelable int `gorm:"not null;default:0" json:"cancelable"`

	// Tracking Available of the Courier Service
	// in: integer
	TrackingAvailable int `gorm:"not null;default:0" json:"tracking_available"`

	// Status of the Courier Service
	// in: integer
	Status *int `gorm:"not null;default:1" json:"status"`

	// Max Weight of the Courier Service
	// in: float64
	MaxWeight float64 `gorm:"not null;default:0" json:"max_weight"`

	// Max Volume of the Courier Service
	// in: float64
	MaxVolume float64 `gorm:"not null;default:0" json:"max_volume"`

	// Max Distance of the Courier Service
	// in: float64
	MaxDistance float64 `gorm:"not null;default:0" json:"max_distance"`

	// Min Purchase of the Courier Service
	// in: float64
	MinPurchase float64 `gorm:"not null;default:0" json:"min_purchase"`

	// Max Purchase of the Courier Service
	// in: float64
	MaxPurchase float64 `gorm:"not null;default:0" json:"max_purchase"`

	// Insurance of the Courier Service
	// in: integer
	Insurance int `gorm:"not null;default:0" json:"insurance"`

	// Insurance Min of the Courier Service
	// in: float64
	InsuranceMin float64 `gorm:"not null;default:0" json:"insurance_min"`

	// Insurance Fee Type of the Courier Service
	// in: string
	InsuranceFeeType string `gorm:"not null" json:"insurance_fee_type"`

	// Insurance Fee of the Courier Service
	// in: float64
	InsuranceFee float64 `gorm:"not null;default:0" json:"insurance_fee"`

	// Start Time of the Courier Service
	// example:"15:04:05+07"
	StartTime datatype.Time `gorm:"null" json:"start_time"`

	// End Time of the Courier Service
	// example:"15:04:05+07"
	EndTime datatype.Time `gorm:"null" json:"end_time"`

	// Repickup Fee of the Courier Service
	// in: float64
	Repickup float64 `gorm:"not null;default:0" json:"repickup"`

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`
}

// swagger:model CourierService
type CourierService struct {
	base.BaseIDModel

	// Courier Id of the Courier Service
	// in: uint64
	CourierID uint64 `gorm:"not null" json:"-"`

	// Courier UId of the Courier Service
	// in: string
	CourierUId string `gorm:"not null" json:"courier_uid"`

	// Shipping Code of the Courier Service
	// in: string
	ShippingCode string `gorm:"not null" json:"shipping_code"`

	// Shipping Name of the Courier Service
	// in: string
	ShippingName string `gorm:"not null" json:"shipping_name"`

	// Shipping Type of the Courier Service
	// in: string
	ShippingType string `gorm:"not null" json:"shipping_type"`

	// Shipping Description of the Courier Service
	// in: string
	ShippingDescription string `json:"shipping_description"`

	// ETD Min of the Courier Service
	// in: float64
	ETD_Min float64 `gorm:"not null" json:"etd_min"`

	// ETD Max of the Courier Service
	// in: float64
	ETD_Max float64 `gorm:"not null" json:"etd_max"`

	// Logo of the Courier Service
	// in: string
	Logo string `gorm:"not null" json:"logo"`

	// Cod Available of the Courier Service
	// in: integer
	CodAvailable int `gorm:"not null;default:0" json:"cod_available"`

	// Prescription Allowed of the Courier Service
	// in: integer
	PrescriptionAllowed int `gorm:"not null;default:0" json:"prescription_allowed"`

	// Cancelable of the Courier Service
	// in: integer
	Cancelable int `gorm:"not null;default:0" json:"cancelable"`

	// Tracking Available of the Courier Service
	// in: integer
	TrackingAvailable int `gorm:"not null;default:0" json:"tracking_available"`

	// Status of the Courier Service
	// in: integer
	Status *int `gorm:"not null;default:1" json:"status"`

	// Max Weight of the Courier Service
	// in: float64
	MaxWeight float64 `gorm:"not null;default:0" json:"max_weight"`

	// Max Volume of the Courier Service
	// in: float64
	MaxVolume float64 `gorm:"not null;default:0" json:"max_volume"`

	// Max Distance of the Courier Service
	// in: float64
	MaxDistance float64 `gorm:"not null;default:0" json:"max_distance"`

	// Min Purchase of the Courier Service
	// in: float64
	MinPurchase float64 `gorm:"not null;default:0" json:"min_purchase"`

	// Max Purchase of the Courier Service
	// in: float64
	MaxPurchase float64 `gorm:"not null;default:0" json:"max_purchase"`

	// Insurance of the Courier Service
	// in: integer
	Insurance int `gorm:"not null;default:0" json:"insurance"`

	// Insurance Min of the Courier Service
	// in: float64
	InsuranceMin float64 `gorm:"not null;default:0" json:"insurance_min"`

	// Insurance Fee Type of the Courier Service
	// in: string
	InsuranceFeeType string `gorm:"not null" json:"insurance_fee_type"`

	// Insurance Fee of the Courier Service
	// in: float64
	InsuranceFee float64 `gorm:"not null;default:0" json:"insurance_fee"`

	// Start Time of the Courier Service
	// in: time
	StartTime datatype.Time `gorm:"null" json:"start_time"`

	// End Time of the Courier Service
	// in: time
	EndTime datatype.Time `gorm:"null" json:"end_time"`

	// Repickup Fee of the Courier Service
	// in: float64
	Repickup float64 `gorm:"not null;default:0" json:"repickup"`

	// Image UID
	// in: string
	ImageUID string `gorm:"size:50;null" json:"image_uid"`

	// Image Path
	// in: string
	// example: [{"path": "image_path", "size": "thumbnail"},{"path": "{image_path}", "size": "original"}]
	ImagePath datatype.JSONB `gorm:"type:jsonb;null" json:"image_path"`

	Courier *Courier `json:"_" gorm:"foreignKey:courier_id"`
}
