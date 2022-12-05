package render

import (
	"fmt"
	"github.com/Go-Sun/sun/sun/internal/bytesconv"
	"net/http"
)

type String struct {
	Format string
	Data   []any
}

func (s *String) Render(w http.ResponseWriter) error {
	s.WriterContentType(w)
	if len(s.Data) > 0 {
		_, err := fmt.Fprintf(w, s.Format, s.Data...)
		return err
	}
	_, err := w.Write(bytesconv.StringToBytes(s.Format))
	return err
}

func (s *String) WriterContentType(w http.ResponseWriter) {
	writeContentType(w, "text/plain; charset=utf-8")
}
