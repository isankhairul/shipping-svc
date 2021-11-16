package entity

//Request Check Device
type JSONRequestProduct struct {
	Name   string `json:"name"`
	Sku    string `json:"sku"`
	Uom    string `json:"uom"`
	Weight int32  `json:"weight"`
}