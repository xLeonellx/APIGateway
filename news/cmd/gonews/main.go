package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/rss"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

const (
	configURL = "./cmd/gonews/config.json"
	dbURL     = "postgres://postgres:postgrespw@db_news:5432/news"
	newsAddr  = ":8080"
)

func main() {

	// инициализация db
	psgr, err := postgres.New(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// инициализация api
	a := api.New(psgr)

	chP := make(chan []storage.Post)
	chE := make(chan error)

	// Чтение RSS-лент из конфига с заданным интервалом
	go func() {
		err := rss.GoNews(configURL, chP, chE)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// запись публикаций в db
	go func() {
		for posts := range chP {
			if err := a.DB.PostsMany(posts); err != nil {
				chE <- err
			}
		}
	}()

	// вывод ошибок
	go func() {
		for err := range chE {
			log.Println(err)
		}
	}()

	// запуск сервера с api
	a.Router().Use(Middle)
	log.Printf("News server is starting on %s", newsAddr)
	err = http.ListenAndServe(newsAddr, a.Router())
	if err != nil {
		log.Fatal("News server could not start. Error:", err)
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
