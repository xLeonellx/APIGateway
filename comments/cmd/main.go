package main

import (
	"comments/pkg/api"
	"comments/pkg/storage/postgres"
	"log"
	"net/http"
)

const (
	dbURL        = "postgres://postgres:postgrespw@db_comments:5432/comments"
	commentsAddr = ":8081"
)

func main() {

	// инициализация db
	psgr, err := postgres.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// инициализация api
	a := api.New(psgr)

	// запуск сервера с api
	a.Router().Use(Middle)
	log.Printf("Comments server is starting on %s", commentsAddr)
	err = http.ListenAndServe(commentsAddr, a.Router())
	if err != nil {
		log.Fatal("Comments server could not start. Error:", err)
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Middle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqID := req.Header.Get("X-Request-ID")

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Printf("<-- client ip: %s, method: %s, url: %s, status code: %d %s, trace id: %s",
			req.RemoteAddr, req.Method, req.URL.Path, statusCode, http.StatusText(statusCode), reqID)

	})
}
