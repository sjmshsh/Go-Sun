package render

import (
	"fmt"
	"log"
	"net/http"
)

type Redirect struct {
	Code     int
	Request  *http.Request
	Location string
}

func (r *Redirect) Render(w http.ResponseWriter) error {
	if (r.Code < http.StatusMultipleChoices ||
		r.Code > http.StatusPermanentRedirect) &&
		http.StatusCreated != r.Code {
		log.Panicln(fmt.Sprintf("Cannot redirect with status code %d", r.Code))
	}
	http.Redirect(w, r.Request, r.Location, r.Code)
	return nil
}

func (r *Redirect) WriterContentType(w http.ResponseWriter) {}
