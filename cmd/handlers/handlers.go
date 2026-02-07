package handlers

import (
	"net/http"
	"time"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Everything is good and running"))
}

func Homehandler(w http.ResponseWriter, r *http.Request) {
	welcome := []string{"Hey", "there", "welcome", "to", "my", "page", "enjoy", "!"}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Write([]byte("Hey there, my name is Diego!\n\n"))
	flusher, _ := w.(http.Flusher)
	for _, s := range welcome {
		w.Write([]byte(s + " "))
		flusher.Flush()
		time.Sleep(time.Millisecond * 500)
	}
}
