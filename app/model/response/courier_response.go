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
