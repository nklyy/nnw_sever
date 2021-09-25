package codes

type Code int

const (
	BadRequest     = 400
	Unauthorized   = 401
	Forbidden      = 403
	NotFound       = 404
	DuplicateError = 409
	InternalError  = 500
)
