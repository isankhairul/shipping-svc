package entity

// swagger:parameters product
type JSONRequestProduct struct {
	// name: Name
	// in: body
	// type: string
	// required: true
	Name string `json:"name"`
	// name: Sku
	// in: body
	// type: string
	// required: true
	Sku string `json:"sku"`
	// name: Uom
	// in: body
	// type: string
	// required: true
	Uom string `json:"uom"`
	// name: Weight
	// in: body
	// type: integer
	// required: true
	Weight int32 `json:"weight"`
}

// swagger:parameters update_product
type JSONRequestUpdateProduct struct {
	// name: id
	// in: path
	// type: integer
	// required: true
	Id string `json:"id"`
	// Name of the Product
	// in: string
	// name: Name
	// in: formData
	// type: string
	// required: true
	Name string `json:"name"`
	// name: Sku
	// in: body
	// type: string
	// required: true
	Sku string `json:"sku"`
	// name: Uom
	// in: body
	// type: string
	// required: true
	Uom string `json:"uom"`
	// name: Weight
	// in: body
	// type: integer
	// required: true
	Weight int32 `json:"weight"`
}

// swagger:parameters getProduct
type getProduct struct {
	// name: id
	// in: path
	// type: string
	// required: true
	Id string `json:"id"`
}

type JSONResponseProduct struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Sku    string `json:"sku"`
	Uom    string `json:"uom"`
	Weight int32  `json:"weight"`
}
