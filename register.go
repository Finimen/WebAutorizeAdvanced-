package main

import "net/http"

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use post", http.StatusBadRequest)
		return
	}
}
