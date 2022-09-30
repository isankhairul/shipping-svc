package entity

import (
	"go-klikdokter/app/model/base"
	"go-klikdokter/app/model/request"
	"go-klikdokter/pkg/util"
	"math"
	"strconv"
	"time"
)

type OrderShipping struct {
	base.BaseIDModel
	OrderNo              string    `gorm:"type:varchar(50);not null"`
	OrderShippingDate    time.Time `gorm:"type:timestamp;not null"`
	ChannelID            uint64    `gorm:"type:bigint;not null"`
	CourierID            uint64    `gorm:"type:bigint;not null"`
	CourierServiceID     uint64    `gorm:"type:bigint;not null"`
	OrderNoAPI           string    `gorm:"type:varchar(50);not null"`
	CustomerUID          string    `gorm:"type:varchar(50);null"`
	CustomerName         string    `gorm:"type:varchar(100);not null"`
	CustomerPhoneNumber  string    `gorm:"type:varchar(20);.not null"`
	CustomerEmail        string    `gorm:"type:varchar(100);null"`
	CustomerAddress      string    `gorm:"type:varchar(255);not null"`
	CustomerLatitude     float64   `gorm:"type:numeric;null"`
	CustomerLongitude    float64   `gorm:"type:numeric;null"`
	CustomerCountryCode  string    `gorm:"type:varchar(50);not null"`
	CustomerProvinceCode string    `gorm:"type:varchar(50);not null"`
	CustomerCityCode     string    `gorm:"type:varchar(50);not null"`
	CustomerDistrictCode string    `gorm:"type:varchar(50);not null"`
	CustomerPostalCode   string    `gorm:"type:varchar(50);not null"`
	CustomerNotes        string    `gorm:"type:varchar(255);null"`
	MerchantUID          string    `gorm:"type:varchar(50);not null"`
	MerchantName         string    `gorm:"type:varchar(100);not null"`
	MerchantPhoneNumber  string    `gorm:"type:varchar(20);not null"`
	MerchantEmail        string    `gorm:"type:varchar(100);null"`
	MerchantAddress      string    `gorm:"type:varchar(255);not null"`
	MerchantLatitude     float64   `gorm:"type:numeric;null"`
	MerchantLongitude    float64   `gorm:"type:numeric;null"`
	MerchantCountryCode  string    `gorm:"type:varchar(50);not null"`
	MerchantProvinceCode string    `gorm:"type:varchar(50);not null"`
	MerchantCityCode     string    `gorm:"type:varchar(50);not null"`
	MerchantDistrictCode string    `gorm:"type:varchar(50);not null"`
	MerchantPostalCode   string    `gorm:"type:varchar(50);not null"`
	TotalWeight          float64   `gorm:"type:numeric;not null"`
	TotalVolume          float64   `gorm:"type:numeric;null"`
	TotalProductPrice    float64   `gorm:"type:numeric;not null"`
	TotalFinalWeight     float64   `gorm:"type:numeric;not null"`
	ContainPrescription  uint      `gorm:"type:numeric;not null"`
	Insurance            bool      `gorm:"type:boolean;null"`
	InsuranceCost        float64   `gorm:"type:numeric;null"`
	ShippingCost         float64   `gorm:"type:numeric;null"`
	TotalShippingCost    float64   `gorm:"type:numeric;null"`
	ActualShippingCost   float64   `gorm:"type:numeric;null"`
	ShippingNotes        string    `gorm:"type:varchar(255);null"`
	BookingID            string    `gorm:"type:varchar(50);null"`
	Airwaybill           string    `gorm:"type:varchar(50);null"`
	Status               string    `gorm:"type:varchar(50);null"`
	PickupCode           string    `gorm:"type:varchar(50);null"`

	Channel              *Channel               `gorm:"foreignKey:channel_id"`
	Courier              *Courier               `gorm:"foreignKey:courier_id"`
	CourierService       *CourierService        `gorm:"foreignKey:courier_service_id"`
	OrderShippingItem    []OrderShippingItem    `gorm:"foreignKey:order_shipping_id"`
	OrderShippingHistory []OrderShippingHistory `gorm:"foreignKey:order_shipping_id"`
}

func (o *OrderShipping) FromCreateDeliveryRequest(req *request.CreateDelivery) {
	cusLat, _ := strconv.ParseFloat(req.Destination.Latitude, 64)
	cusLong, _ := strconv.ParseFloat(req.Destination.Longitude, 64)
	merLat, _ := strconv.ParseFloat(req.Origin.Latitude, 64)
	merLong, _ := strconv.ParseFloat(req.Origin.Longitude, 64)
	volumeWeight := util.CalculateVolumeWeightKg(req.Package.TotalLength, req.Package.TotalWeight, req.Package.TotalHeight)

	var orderShippingItems []OrderShippingItem

	for _, v := range req.Package.Product {
		orderShippingItems = append(orderShippingItems, OrderShippingItem{
			ItemName:   v.Name,
			ProductUID: v.UID,
			Price:      v.Price,
			Quantity:   v.Qty,
			TotalPrice: v.Price * float64(v.Qty),
			BaseIDModel: base.BaseIDModel{
				CreatedBy: req.ActorName,
			},
		})
	}

	o.OrderNo = req.OrderNo
	o.OrderShippingDate = time.Now()
	o.ChannelID = 0
	o.CourierID = 0
	o.CourierServiceID = 0
	o.OrderNoAPI = req.OrderNo
	o.CustomerUID = req.Customer.UID
	o.CustomerName = req.Customer.Name
	o.CustomerPhoneNumber = req.Customer.Phone
	o.CustomerEmail = req.Customer.Email
	o.CustomerAddress = req.Destination.Address
	o.CustomerLatitude = cusLat
	o.CustomerLongitude = cusLong
	o.CustomerCountryCode = req.Destination.CountryCode
	o.CustomerProvinceCode = req.Destination.ProvinceCode
	o.CustomerCityCode = req.Destination.CityName
	o.CustomerDistrictCode = req.Destination.DistrictName
	o.CustomerPostalCode = req.Destination.PostalCode
	o.CustomerNotes = req.Notes
	o.MerchantUID = req.Merchant.UID
	o.MerchantName = req.Merchant.Name
	o.MerchantPhoneNumber = req.Merchant.Phone
	o.MerchantEmail = req.Merchant.Email
	o.MerchantAddress = req.Origin.Address
	o.MerchantLatitude = merLat
	o.MerchantLongitude = merLong
	o.MerchantCountryCode = req.Destination.CountryCode
	o.MerchantProvinceCode = req.Destination.ProvinceCode
	o.MerchantCityCode = req.Destination.CityName
	o.MerchantDistrictCode = req.Destination.DistrictName
	o.MerchantPostalCode = req.Destination.PostalCode
	o.TotalWeight = req.Package.TotalWeight
	o.TotalVolume = util.CalculateVolume(req.Package.TotalLength, req.Package.TotalWeight, req.Package.TotalHeight)
	o.TotalProductPrice = req.Package.TotalProductPrice
	o.TotalFinalWeight = math.Max(volumeWeight, req.Package.TotalWeight)
	o.ContainPrescription = req.Package.ContainPrescription
	o.ShippingNotes = req.Notes
	o.OrderShippingItem = orderShippingItems
	o.OrderShippingHistory = []OrderShippingHistory{}
	o.BaseIDModel = base.BaseIDModel{
		CreatedBy: req.ActorName,
		UpdatedBy: req.ActorName,
	}
}
func (o *OrderShipping) AddHistoryStatus(s *ShippingCourierStatus, note string) {
	if o.isHistoryStatusExist(s.StatusCode, note) {
		return
	}

	o.OrderShippingHistory = append(o.OrderShippingHistory, OrderShippingHistory{
		OrderShippingID:         o.ID,
		ShippingCourierStatusID: s.ID,
		StatusCode:              s.StatusCode,
		Note:                    note,
		BaseIDModel: base.BaseIDModel{
			CreatedBy: o.UpdatedBy,
		},
	})
}

func (o *OrderShipping) isHistoryStatusExist(statusCode, note string) bool {
	for _, v := range o.OrderShippingHistory {
		if v.StatusCode == statusCode && v.Note == note {
			return true
		}
	}

	return false
}

func (OrderShipping) TableName() string {
	return "order_shipping"
}
