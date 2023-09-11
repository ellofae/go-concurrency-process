package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ellofae/go-concurrency-process/pkg/logger"
)

func ToJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(i)
}

func StructDecode(r *http.Request, req interface{}) error {
	logger := logger.GetLogger()

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		logger.Error("Unable to decode the request data", "error", err)
		return err
	}

	return nil
}
