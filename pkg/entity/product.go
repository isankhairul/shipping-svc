package entity

type Product struct {
	BaseIDModel
	Name   string `json:"name"`
	Sku    string `json:"sku"`
	Uom    string `json:"uom"`
	Weight int32  `json:"weight"`
}
