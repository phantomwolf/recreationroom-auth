package response

import (
	"encoding/json"
	"github.com/volatiletech/null"
)

// Response is an all-purpose type for JSON responses
// Stats: "ok" or "error"
// Code: status code which indicates the type of error
// Messages: an array of error messages
// Result: payload, can be anything
// Err: error
type Response struct {
	Status   string                 `json:"status"`
	Code     int                    `json:"code"`
	Result   map[string]interface{} `json:"result"`
	Messages []string               `json:"messages"`
	Err      error                  `json:"-"`
}

func New(status string, code int, err error, msgs ...string) *Response {
	return &Response{
		Status:   status,
		Code:     code,
		Result:   map[string]interface{}{},
		Messages: msgs,
		Err:      err,
	}
}

func (res *Response) SetStatus(status string, code int, err error) {
	res.Status = status
	res.Code = code
	res.Err = err
}

func (res *Response) SetResult(key string, v interface{}) {
	res.Result[key] = v
}

func (res *Response) AddMessage(msg ...string) {
	res.Messages = append(res.Messages, msg...)
}

func (res *Response) SetError(err error) {
	res.Err = err
}

func (res *Response) MarshalJSON() ([]byte, error) {
	type Alias Response
	var err null.String
	if res.Err != nil {
		err = null.StringFrom(res.Err.Error())
	}
	return json.Marshal(&struct {
		*Alias
		Err null.String `json:"error"`
	}{
		Alias: (*Alias)(res),
		Err:   err,
	})
}
