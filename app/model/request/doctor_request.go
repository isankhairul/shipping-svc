package request

// swagger:model SaveProductRequest
type SaveDoctorRequest struct {
	// Name of the product
	// in: string
	Name string `json:"name"`

	// Sku of the product
	// in: string
	Gender string `json:"gender"`

	// Uid of the product, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`
}
