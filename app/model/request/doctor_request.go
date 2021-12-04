package request

// swagger:parameters SaveDoctorRequest
type ReqDoctorBody struct {
	//  in: body
	Body SaveDoctorRequest `json:"body"`
}

type SaveDoctorRequest struct {
	// Name of the doctor
	// in: string
	Name string `json:"name"`

	// Gender of the doctor
	// in: string
	Gender string `json:"gender"`

	// Uid of the product, use this on UPDATE function
	// in: int32
	Uid string `json:"uid" binding:"omitempty"`
}
