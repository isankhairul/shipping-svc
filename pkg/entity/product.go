package entity

// swagger:model Company
type Product struct {
	BaseIDModel
	// Name of the product
	// in: string
	Name string `json:"name"`
	// Sku of the product
	// in: string
	Sku string `json:"sku"`
	// Uom of the product
	// in: string
	Uom string `json:"uom"`
	// Weight of the product
	// in: int32
	Weight int32 `json:"weight"`
}
