package Tool

import (
	"encoding/json"
	"io"
)

//解析json
func Decode(io io.ReadCloser, v interface{}) error{
	return json.NewDecoder(io).Decode(v)
}
