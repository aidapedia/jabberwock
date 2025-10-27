package log

import (
	"net/http"

	cerror "github.com/aidapedia/jabberwock/internal/common/error"
)

func Metadata(code int, message string) map[string]interface{} {
	if message == "" {
		message = cerror.ErrorMessageTryAgain.Error()
	}
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return map[string]interface{}{
		"server_code":    code,
		"server_message": message,
	}
}
