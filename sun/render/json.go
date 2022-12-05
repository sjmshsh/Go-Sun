package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

func (r *JSON) Render(w http.ResponseWriter) error {
	r.WriterContentType(w)
	jsonBytes, err := json.Marshal(r.Data)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func (r *JSON) WriterContentType(w http.ResponseWriter) {
	writeContentType(w, "application/json; charset=utf-8")

}
