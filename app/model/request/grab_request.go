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

	// for create order
	Keywords string `json:"keywords,omitempty"`
}

type Destination struct {
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`

	// for create order
	Keywords string `json:"keywords,omitempty"`
}

type GrabDeliveryQuotes struct {
	ServiceType string      `json:"serviceType"`
	Origin      Origin      `json:"origin"`
	Destination Destination `json:"destination"`
	Packages    []Package   `json:"packages"`
}

type CreateDeliveryGrab struct {
	MerchantOrderID string    `json:"merchantOrderID"`
	ServiceType     string    `json:"serviceType"`
	PaymentMethod   string    `json:"paymentMethod"`
	Packages        []Package `json:"packages"`
	//CashOnDelivery  CashOnDelivery      `json:"cashOnDelivery"`
	Sender      GrabSenderRecipient `json:"sender"`
	Recipient   GrabSenderRecipient `json:"recipient"`
	Origin      Origin              `json:"origin"`
	Destination Destination         `json:"destination"`
	Schedule    Schedule            `json:"schedule"`
}

type CashOnDelivery struct {
	Amount int `json:"amount"`
}
type GrabSenderRecipient struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Title       string `json:"title"`
	CompanyName string `json:"companyName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	SmsEnabled  bool   `json:"smsEnabled"`
	Instruction string `json:"instruction"`
}
type Schedule struct {
	// Date/Time IETF RFC 3339
	PickupTimeFrom string `json:"pickupTimeFrom"`

	// Date/Time IETF RFC 3339
	PickupTimeTo string `json:"pickupTimeTo"`
}

type Extra struct {
}

type GrabAuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
}
