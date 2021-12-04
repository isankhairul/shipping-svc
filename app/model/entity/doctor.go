package entity

import "go-klikdokter/app/model/base"

// swagger:model doctor
type Doctor struct {
	base.BaseIDModel
	// Name of the doctor
	// in: string
	Name string `json:"name"`

	// Gender of the doctor
	// in: string
	Gender string `json:"gender"`
}
