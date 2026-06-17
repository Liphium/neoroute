package neoroute

//go:generate msgp -unexported

type response struct {
	Id      int    `msg:"id"`
	HasData bool   `msg:"has_data"`
	IsError bool   `msg:"error"`
	Data    []byte `msg:"data"`
}

type responseData struct {
	HasData bool   `msg:"has_data"`
	IsError bool   `msg:"error"`
	Data    []byte `msg:"data"`
}

// NewError creates a new error response with the given message.
// This will sent the error message as an error response to the user.
func NewError(msg string) error {
	return &responseData{
		HasData: true,
		IsError: true,
		Data:    []byte(msg),
	}
}

func (r *responseData) Error() string {
	return ""
}

func (r *responseData) Is(target error) bool {
	_, ok := target.(*responseData)
	return ok
}

// This type is used for routes that have no response so no error is thrown.
type noResponse struct{}

func (r *noResponse) Error() string {
	return ""
}

func (r *noResponse) Is(target error) bool {
	_, ok := target.(*noResponse)
	return ok
}
