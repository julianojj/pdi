package exceptions

type ValidationException struct {
	Message string
}

func NewValidationException(message string) error {
	return &ValidationException{
		Message: message,
	}
}

func (e *ValidationException) Error() string {
	return e.Message
}
