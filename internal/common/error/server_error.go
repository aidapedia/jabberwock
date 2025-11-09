package error

import (
	"database/sql"
	"errors"
)

var (
	ErrorMessageTryAgain = errors.New("Oops, something went wrong. Please try again later.")
)

// DB Error
var (
	ErrorNotFound = sql.ErrNoRows
)
