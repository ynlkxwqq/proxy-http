package response

import (
	"net/http"

	"github.com/go-chi/render"
)

type Object struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func OK(w http.ResponseWriter, r *http.Request, data any) {
	render.Status(r, http.StatusOK)

	v := Object{
		Message: "success",
		Data:    data,
	}

	render.JSON(w, r, v)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)

	v := Object{
		Message: "not found",
	}

	render.JSON(w, r, v)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)

	v := Object{
		Message: err.Error(),
	}

	render.JSON(w, r, v)
}
