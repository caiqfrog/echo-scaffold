package utils

import (
	"io"
	"encoding/json"
)

func DecodeArgs(r io.Reader, ptr interface{}) error {
	decoder := json.NewDecoder(r)
	if e := decoder.Decode(ptr); nil != e && io.EOF != e {
		return e
	}
	return nil
}
