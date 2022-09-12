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

var WalletNotFoundError = ErrorResponse{
	Error: "wallet not found",
}

func ErrWalletNotFound(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &WalletNotFoundError)
}

var WalletLockedError = ErrorResponse{
	Error: "wallet locked",
}

func ErrWalletLocked(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &WalletLockedError)
}

var WalletNotLockedError = ErrorResponse{
	Error: "wallet not locked",
}

func ErrWalletNotLocked(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &WalletNotLockedError)
}

var InvalidKeyError = ErrorResponse{
	Error: "Invalid key",
}

func ErrInvalidKey(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &InvalidKeyError)
}

var WalletNoPasswordError = ErrorResponse{
	Error: "password not set",
}

func ErrNoWalletPassword(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &WalletNoPasswordError)
}

var InvalidHashError = ErrorResponse{
	Error: "Invalid hash",
}

func ErrInvalidHash(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, &InvalidHashError)
}

var WorkFailedError = ErrorResponse{
	Error: "Failed to generate work",
}

func ErrWorkFailed(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, &WorkFailedError)
}

func ErrInternalServerError(w http.ResponseWriter, r *http.Request, errorText string) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, &ErrorResponse{
		Error: errorText,
	})
}
