package request

type Dimensions struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Depth  int `json:"depth"`
	Weight int `json:"weight"`
}

type Package struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Quantity    int        `json:"quantity"`
	Price       int        `json:"price"`
	Dimensions  Dimensions `json:"dimensions"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Origin struct {
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}

type Destination struct {
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
}

type GrabDeliveryQuotes struct {
	ServiceType string      `json:"serviceType"`
	Origin      Origin      `json:"origin"`
	Destination Destination `json:"destination"`
	Packages    []Package   `json:"packages"`
}
