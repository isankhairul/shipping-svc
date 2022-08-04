package response

import "go-klikdokter/app/model/entity"

// swagger:model ShippingTypeItem
type ShippingTypeItem struct {

	//example: pqC8LqdgT2KKdcmn2bHnR
	UID string `json:"uid"`

	//example: same_day
	ShippingTypeCode string `json:"shipping_type_code"`

	//example: Same Day Delivery
	Name string `json:"name"`
}

// swagger:response ShippingTypeList
type ShippingTypeList struct {
	//in: body
	ResponseBody []ShippingTypeItem `json:"response"`
}

func NewShippingTypeItem(input entity.ShippmentPredefined) ShippingTypeItem {
	return ShippingTypeItem{
		UID:              input.UID,
		ShippingTypeCode: input.Code,
		Name:             input.Title,
	}
}

func NewShippingTypeItemList(input []entity.ShippmentPredefined) []ShippingTypeItem {
	result := []ShippingTypeItem{}

	for _, v := range input {
		result = append(result, NewShippingTypeItem(v))
	}

	return result
}

//swagger:model CourierByChannelResponse
type CourierByChannelResponse struct {
	//example:aabb7778888dddeeeee
	CourierUID string `json:"courier_uid"`

	//example:shipper
	CourierCode string `json:"courier_code"`

	//example:Shipper
	CourierName string `json:"courier_name"`

	//example:third_party
	CourierTypeCode string `json:"courier_type_code"`

	//example:Third Party Courier
	CourierTypeName string `json:"courier_type_name"`

	//example:https://kd-test/logo/shipper.logo
	ImageLogo string `json:"image_logo" gorm:"courier_image"`
}

//swagger:model CourierServiceByChannel
type CourierServiceByChannelResponse struct {
	//example:ssssbb7778888dddzzzzz
	CourierServiceUID string `json:"courier_service_uid"`

	//example:shipper
	ShippingCode string `json:"shipping_code"`

	//example:Shipper
	ShippingName string `json:"shipping_name"`

	//example:Shipper adalah paket reguler yang ditawarkan Shipper
	ShippingDescription string `json:"shipping_description"`

	//example:https://kd-test/logo/shipper.logo
	ImageLogo string `json:"image_logo"`

	//example:regular
	ShippingTypeCode string `json:"shipping_type_code"`

	//example:Regular Delivery
	ShippingTypeName string `json:"shipping_type_name"`

	//example:24
	ETDMin float64 `json:"etd_min"`

	//example:48
	ETDMax float64 `json:"etd_max"`

	Courier CourierByChannelResponse `json:"courier" gorm:"-:all"`

	CourierUID      string `json:"-"`
	CourierCode     string `json:"-"`
	CourierName     string `json:"-"`
	CourierTypeCode string `json:"-"`
	CourierTypeName string `json:"-"`
	CourierImage    string `json:"-"`
}

//swagger:response CourierByChannel
type CourierByChannelResponseList struct {
	//in: body
	Response []CourierServiceByChannelResponse `json:"response"`
}
