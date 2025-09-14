package log

func Metadata(code int, message string) map[string]interface{} {
	return map[string]interface{}{
		"server_code":    code,
		"server_message": message,
	}
}
