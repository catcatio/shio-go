package endpoints

import (
	"fmt"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintf(w, "err: %s", msg)
}

func writeOKResponse(w http.ResponseWriter) {
	writeResponse(w, http.StatusOK)
}

func writeBadRequest(w http.ResponseWriter) {
	writeResponse(w, http.StatusBadRequest)
}

func writeNotFoundResponse(w http.ResponseWriter) {
	writeResponse(w, http.StatusNotFound)
}

func writeResponse(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	_, _ = fmt.Fprintf(w, http.StatusText(status))
}
