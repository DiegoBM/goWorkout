package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Envelope map[string]interface{}

func WtiteJSON(w http.ResponseWriter, status int, data Envelope) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ReadIDParam(r *http.Request) (int64, error) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		return -1, errors.New("invalid param \"id\"")
	}

	id, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		return -1, errors.New("invalid param type for \"id\"")
	}

	return id, nil
}
