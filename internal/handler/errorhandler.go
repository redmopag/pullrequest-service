package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/avito/internship/pr-service/internal/model"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func ErrorHandler(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			renderError(w, err)
		}
	}
}

var errorStatusMap = map[string]int{
	"NOT_FOUND":    http.StatusNotFound,
	"TEAM_EXISTS":  http.StatusConflict,
	"PR_EXISTS":    http.StatusConflict,
	"PR_MERGED":    http.StatusConflict,
	"NOT_ASSIGNED": http.StatusConflict,
	"NO_CANDIDATE": http.StatusConflict,
	"BAD_REQUEST":  http.StatusBadRequest,
}

func renderError(w http.ResponseWriter, err error) {
	log.Printf("unhandled error: %v", err)

	var domainErr *model.DomainError
	if errors.As(err, &domainErr) {
		if status, ok := errorStatusMap[domainErr.Code]; ok {
			WriteJSONResponse(w, status, domainErr)
			return
		}
	}
	WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
