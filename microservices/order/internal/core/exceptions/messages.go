package exceptions

var (
	ErrUserIDIsRequired       = NewValidationException("user id is required")
	ErrPaymentTokenIsRequired = NewValidationException("payment token is required")
	ErrInvalidPaymentToken    = NewValidationException("invalid payment token")
	ErrItemNotFound           = NewNotFoundException("item not found")
	ErrOrderNotFound          = NewNotFoundException("order not found")
)
