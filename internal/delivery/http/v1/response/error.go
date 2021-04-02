package response

type ErrorResponseDto struct {
	Code        int      `json:"code"`
	Message     string   `json:"message"`
	Information []string `json:"information"`
}

func NewErrorResponseDto(code int, message string, information ...string) ErrorResponseDto {
	return ErrorResponseDto{
		Code:        code,
		Message:     message,
		Information: information,
	}
}
