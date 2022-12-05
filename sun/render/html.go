package render

import (
	"github.com/Go-Sun/sun/sun/internal/bytesconv"
	"html/template"
	"net/http"
)

type HTMLRender struct {
	Template *template.Template
}

type HTML struct {
	Data       any
	Name       string
	Template   *template.Template
	IsTemplate bool
}

func (h *HTML) Render(w http.ResponseWriter) error {
	h.WriterContentType(w)
	if h.IsTemplate {
		err := h.Template.ExecuteTemplate(w, h.Name, h.Data)
		return err
	}
	// 如果不是模板的话那么我们规定这里必须传的是string类型
	_, err := w.Write(bytesconv.StringToBytes(h.Data.(string)))
	return err
}

func (h *HTML) WriterContentType(w http.ResponseWriter) {
	writeContentType(w, "text/html; charset=utf-8")
}
