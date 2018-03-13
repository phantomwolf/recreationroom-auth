package response

// Response is an all-purpose type for JSON responses
// Stats: "ok" or "error"
// Code: status code which indicates the type of error
// Messages: an array of error messages
// Result: payload, can be anything
type JSONResponse struct {
	Status   string      `json:"status"`
	Code     int         `json:"code"`
	Messages []string    `json:"messages"`
	Result   interface{} `json:"result"`
}

func NewJSONResponse(status string, code int, messages []string, result interface{}) *JSONResponse {
	return &JSONResponse{
		Status:   status,
		Code:     code,
		Messages: messages,
		Result:   result,
	}
}
