package error

import "errors"

var (
	ErrorMessageTryAgain = errors.New("Oops, something went wrong. Please try again later.")
)

// DB Error
var (
	ErrorNotFound = errors.New("not found")
)
