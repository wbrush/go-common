package errorhandler

const (
	ErrService       ErrorCode = "ERR_SERVICE"
	ErrNotFound      ErrorCode = "ERR_NOT_FOUND"
	ErrAlreadyExists ErrorCode = "ERR_ALREADY_EXISTS"
	ErrBadRequest    ErrorCode = "ERR_BAD_REQUEST"
	ErrBadParam      ErrorCode = "ERR_BAD_PARAM"
	ErrNotAllowed    ErrorCode = "ERR_NOT_ALLOWED"
	ErrUnauthorized  ErrorCode = "ERR_UNAUTHORIZED"
)

func (ec ErrorCode) IsValid() bool {
	switch ec {
	case ErrService, ErrNotFound, ErrAlreadyExists, ErrBadRequest,
		ErrBadParam, ErrNotAllowed, ErrUnauthorized:
		return true
	default:
		return false
	}
}
