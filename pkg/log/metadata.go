package log

import (
	"net/http"

	"github.com/aidapedia/jabberwock/internal/common"
)

func Metadata(code int, message string) map[string]interface{} {
	if message == "" {
		message = common.ErrorMessageTryAgain
	}
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return map[string]interface{}{
		"server_code":    code,
		"server_message": message,
	}
}
