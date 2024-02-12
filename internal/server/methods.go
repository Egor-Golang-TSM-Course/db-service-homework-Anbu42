package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) respond(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data = map[string]any{
		"response": data,
	}

	return json.NewEncoder(w).Encode(data)
}
