package entity

import "go-klikdokter/app/model/base"

// swagger:model product
type Doctor struct {
	base.BaseIDModel
	// Name of the product
	// in: string
	Name string `json:"name"`

	// Uom of the product
	// in: string
	Gender string `json:"gender"`
}
