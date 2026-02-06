package server

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}
