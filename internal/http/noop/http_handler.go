package noop

import (
	"log"
	"net/http"
)

// Handler noop HTTP Hander. Does nothing, returns HTTP 200 OK
func Handler(w http.ResponseWriter, r *http.Request) {
	log.Println("[DEBUG] noop")
}
