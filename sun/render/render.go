package render

import "net/http"

type Render interface {
	Render(w http.ResponseWriter) error
	WriterContentType(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, value string) {
	w.Header().Set("Content-type", value)
}
