package entity

import (
	"go-klikdokter/app/model/base"
)

// swagger:model product
type Product struct {
	base.BaseIDModel
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
