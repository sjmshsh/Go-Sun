package sun

import (
	"log"
	"net/http"
)

type Engine struct {

}

func New () *Engine {
	return &Engine{}
}

func (e *Engine) Run() {
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		log.Fatal(err)
	}
}
