package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data   any
}

func (r *XML) Render(w http.ResponseWriter) error {
	r.WriterContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

func (r *XML) WriterContentType(w http.ResponseWriter) {
	writeContentType(w, "application/xml; charset=utf-8")
}
