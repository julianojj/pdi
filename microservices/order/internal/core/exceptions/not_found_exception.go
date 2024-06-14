package exceptions

type NotFoundException struct {
	Message string
}

func NewNotFoundException(message string) error {
	return &NotFoundException{
		Message: message,
	}
}

func (e *NotFoundException) Error() string {
	return e.Message
}
