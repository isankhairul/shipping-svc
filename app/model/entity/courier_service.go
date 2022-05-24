package entity

import (
	"go-klikdokter/app/model/base"
	"time"
)

// swagger:model Courier Service
type CourierService struct {
	base.BaseIDModel
	// Courier UId of the Courier Service
	// in: int
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
	Status int `gorm:"not null;default:1" json:"status"`

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
	StartTime time.Time `gorm:"not null" json:"start_time"`

	// End Time of the Courier Service
	// in: time
	EndTime time.Time `gorm:"not null" json:"end_time"`

	// Repickup Fee of the Courier Service
	// in: float64
	Repickup float64 `gorm:"not null;default:0" json:"repickup"`
}
