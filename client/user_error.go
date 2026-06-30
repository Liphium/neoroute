package client

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
