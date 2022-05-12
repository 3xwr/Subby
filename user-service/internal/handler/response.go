package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(b))
}

func getUserID(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	tokens := strings.SplitN(authHeader, " ", 2)
	if len(tokens) != 2 {
		return "", fmt.Errorf("incorrect user ID")
	}
	return tokens[1], nil
}
