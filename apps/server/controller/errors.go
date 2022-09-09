package controller

import (
	"net/http"

	"github.com/go-chi/render"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var UnableToParseJsonError = ErrorResponse{
	Error: "Unable to parse json",
}

func ErrUnableToParseJson(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &UnableToParseJsonError)
}

var InvalidSeedError = ErrorResponse{
	Error: "Invalid seed",
}

func ErrInvalidSeed(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &InvalidSeedError)
}

func ErrInternalServerError(w http.ResponseWriter, r *http.Request, errorText string) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, &ErrorResponse{
		Error: errorText,
	})
}