package tools

import (
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
)

func RespondWithCachedJSON(w http.ResponseWriter, cache string, v any, code int) {
	err := json.Unmarshal([]byte(cache), &v)
	if err != nil {
		RespondWithError(w, errors.New("failed to unmarshall cached book"), http.StatusInternalServerError)
		return
	}

	RespondWithJSON(w, v, code)
}

func RespondWithJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		RespondWithError(w, err, http.StatusInternalServerError)
	}
}

func RespondWithJSONAndCache(w http.ResponseWriter, r *http.Request, client *redis.Client, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	cachedJSON, err := json.Marshal(v)
	if err != nil {
		RespondWithError(w, errors.New("failed to process data"), http.StatusInternalServerError)
		return
	}

	cachedKey := r.Context().Value("cachedKey").(string)
	err = client.SetEx(r.Context(), cachedKey, cachedJSON, time.Second*3600).Err()
	if err != nil {
		log.Printf("failed to cache data: %s", err)
	}

	_, err = w.Write(cachedJSON)
	if err != nil {
		RespondWithError(w, err, http.StatusInternalServerError)
		return
	}
}

func RespondWithError(w http.ResponseWriter, err error, code int) {
	RespondWithJSON(w, map[string]string{"error": err.Error()}, code)
}
