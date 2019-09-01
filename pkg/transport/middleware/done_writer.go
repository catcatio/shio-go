package middleware

import "net/http"

type doneWriter struct {
	http.ResponseWriter
	done bool
}

func (w *doneWriter) WriteHeader(status int) {
	w.done = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *doneWriter) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}

func DoneWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dw := &doneWriter{ResponseWriter: w}
		next.ServeHTTP(dw, r)
		if dw.done {
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
