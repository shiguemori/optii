package utils

import "encoding/json"

// Response Error message response
// swagger:model utils.Response
type Response struct {
	Message string `json:"message"`
}

func ToJSON(o interface{}) string {
	bytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
