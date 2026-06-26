package neogen

type Schema struct {
	Version      int                    `json:"version"`
	Generator    string                 `json:"generator"`
	Transporters map[string]Transporter `json:"transporters"`
}

const (
	TransporterHTTP = iota
	TransporterWebTransport
	TransporterWebSocket
)

type Transporter struct {
	Type   int                  `json:"type"`
	Events map[string]BasicType `json:"events"`
	Routes map[string]Route     `json:"routes"`
}

type Route struct {
	HasRequest bool      `json:"has_request"`
	Request    BasicType `json:"request"`

	HasResponse    bool      `json:"has_response"`
	CanReturnError bool      `json:"can_return_error"`
	Response       BasicType `json:"response"`
}
