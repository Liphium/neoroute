package client

// UserError is an error returned by the server intended for the user.
type UserError struct {
	message string
}

func NewUserError(message string) *UserError {
	return &UserError{
		message: message,
	}
}

func (e *UserError) Error() string {
	return e.message
}
