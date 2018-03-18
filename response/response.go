package response

import (
	"encoding/json"
)

// Response is an all-purpose type for JSON responses
// Stats: "ok" or "error"
// Code: status code which indicates the type of error
// Messages: an array of error messages
// Result: payload, can be anything
type Response struct {
	Status   string                 `json:"status"`
	Code     int                    `json:"code"`
	Result   map[string]interface{} `json:"result"`
	Messages []string               `json:"messages"`
	Err      rror                   `json:"error"`
}

func New() *Response {
	return &Response{
		Status:   "ok",
		Code:     0,
		Result:   map[string]interface{}{},
		Messages: []string{},
		Err:      nil,
	}
}

func (res *Response) SetStatus(status string, code int) {
	res.Status = status
	res.Code = code
}

func (res *Response) SetResult(key string, v interface{}) {
	res.Result[key] = v
}

func (res *Response) AddMessage(msg ...string) {
	res.Messages = append(res.Messages, msg...)
}

func (res *Response) AddError(err ...error) {
	res.Errs = append(res.Errs, err...)
}

func (res *Response) MarshalJSON() ([]byte, error) {
	errs := []string{}
	for _, e := range res.Errs {
		errs = append(errs, e.Error())
	}

	type Alias Response
	return json.Marshal(&struct {
		*Alias
		Errs []string `json:"errors"`
	}{
		Alias: (*Alias)(res),
		Errs:  errs,
	})
}
