package lib

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

// gets status code based on error returned
func GetStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	log.Error(err)
	switch err {
	case ErrInternalServerError:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
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

// TODO: remove interface shit and use composition
func GetJSON(r *http.Request, value interface{}) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&value); err != nil {
		msg := err.Error()
		if serr, ok := err.(*json.SyntaxError); ok {
			log.Error(serr.Error())
			msg += ", at offset: " + strconv.FormatInt(serr.Offset, 10)
		}

		return errors.New("Couldn't process JSON payload, Error: " + msg)
	}

	return nil
}
