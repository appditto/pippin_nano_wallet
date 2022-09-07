package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

type HttpController struct {
}

func (hc *HttpController) HandleAction(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{"hello": "world"})
}
