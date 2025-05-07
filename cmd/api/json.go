package main

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	TotalData   int64    `json:"total_data"`
	SuccessData int64    `json:"success_data"`
	FailData    int64    `json:"fail_data"`
	Error       []string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func writeJSONErrors(w http.ResponseWriter, status int, message []string) error {
	type envelope struct {
		Error []string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}

// func writeJSONResponse(w http.ResponseWriter, status int) error {
// 	var res Response
// 	return writeJSON(w, status, res)
// }
