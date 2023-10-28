package httpserver

import (
	"encoding/json"
	"net/http"
)

type StatusResponse struct {
	Running bool `json:"running"`
}

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(StatusResponse{Running: true})
}
