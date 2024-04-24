package httpserver

import (
	"encoding/json"
	"net/http"
	"net/textproto"
)

// Error represents an error for an end user.
type Error struct {
	Msg string `json:"message,omitempty"`
}

type Response struct {
	Msg string `json:"message,omitempty"`
}

// JSON responds with an error and status code from the error.
func JSON(w http.ResponseWriter, resp any, code int) {
	if err := sendJSON(w, code, resp, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func File(w http.ResponseWriter, resp []byte, code int, headers http.Header) {
	if err := sendFile(w, code, resp, headers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// sendJSON sends a JSON response with a given status.
// In case of an error, response (and status) is not send and error is returned.
func sendJSON(w http.ResponseWriter, status int, resp any, headers http.Header) error {
	js, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	for k, v := range headers {
		k = textproto.CanonicalMIMEHeaderKey(k)
		w.Header()[k] = v
	}

	w.WriteHeader(status)
	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func sendFile(w http.ResponseWriter, status int, resp []byte, headers http.Header) error {
	for k, v := range headers {
		k = textproto.CanonicalMIMEHeaderKey(k)
		w.Header()[k] = v
	}

	w.WriteHeader(status)
	_, err := w.Write(resp)
	if err != nil {
		return err
	}

	return nil
}
