package dto

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

type ListResponseFields struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

type PaginationFields struct {
	Count  int64 `json:"count"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

func SuccessJSONResponse(w http.ResponseWriter, r *http.Request, object interface{}) {
	render.JSON(w, r, object)
}

func SuccessJSON(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, Response{
		StatusCode: http.StatusOK,
		Message:    message,
	})
}

func BadRequestJSON(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, Response{
		StatusCode: http.StatusBadRequest,
		Message:    fmt.Sprintf("Bad Request : %v", message),
	})
}

func ConflictErrorJSON(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusConflict)
	render.JSON(w, r, Response{
		StatusCode: http.StatusConflict,
		Message:    fmt.Sprintf("Conflict Error : %v", message),
	})
}

func InternalServerErrorJSON(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, Response{
		StatusCode: http.StatusInternalServerError,
		Message:    fmt.Sprintf("Internal Server Error : %v", message),
	})
}

func NotFoundJSON(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, Response{
		StatusCode: http.StatusNotFound,
		Message:    fmt.Sprintf("Entity Not Found : %v", message),
	})
}
