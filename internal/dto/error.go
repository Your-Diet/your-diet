package dto

type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewError(field, message string) *Error {
	return &Error{
		Field:   field,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}
