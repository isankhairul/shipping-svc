package response

type ImportStatus struct {
	UID        string `json:"uid,omitempty"`
	CourierUID string `json:"courier_uid"`
	Status     bool   `json:"status"`
	Message    string `json:"message"`
}
