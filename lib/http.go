package lib

import (
	"encoding/json"
	"go_boilerplate/models"
	"net/http"
)

type HttpMessage struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

// respondJSON makes the response with payload as json format
func RespondJSON(w http.ResponseWriter, status int, payload interface{}, error string) {
	response, err := json.Marshal(HttpMessage{Data: payload, Error: error})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// gets status code based on model err
func GetStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	log.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
func UseString(s *string) string {
	if s == nil {
		temp := "" // *string cannot be initialized
		s = &temp  // in one statement
	}
	value := *s // safe to dereference the *string
	return value
}
