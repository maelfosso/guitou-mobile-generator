package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"guitou.cm/mobile/generator/models"
)

func decodeGenerateRequest(_ context.Context, r *http.Request, pam map[string]string) (interface{}, error) {
	if id, ok := pam["id"]; ok {
		return serviceRequest{
			ProjectID: id,
		}, nil
	}

	return nil, errors.New("No ProjectID found")
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if e, ok := response.(error); ok && e == ErrorNoProjectID {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.Error(),
		})
		return nil
	}

	if e, ok := response.(models.ErrorProjectRepository); ok && e.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.Error(),
		})
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(response interface{}, w http.ResponseWriter) {

	if e, ok := response.(error); ok && e == ErrorNoProjectID {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": e.Error(),
		})
	} else {
		if e, ok := response.(models.ErrorProjectRepository); ok && e.Err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			json.NewEncoder(w).Encode(map[string]interface{}{
				"error": e.Error(),
			})
		}
	}
}
