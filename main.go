package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/boramalper/url2id/pkg"
)

func main() {
	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := &http.Server{
		Addr:           ":" + port,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Serving on port %s", port)
	log.Fatal(s.ListenAndServe())
}

func handler(w http.ResponseWriter, r *http.Request) {
	urls, ok := r.URL.Query()["url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("url parameter is missing. see https://github.com/boramalper/url2id for documentation."))
		return
	}

	url_, err := url.Parse(urls[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("url cannot be parsed"))
		return
	}

	result, err2 := pkg.URL2ID(url_)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		if err2.PrevErr != nil {
			log.Printf("URL2ID Error: %s <%s>", err2.PrevErr.Error(), url_.String())
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=21600")
	_ = json.NewEncoder(w).Encode(result)
}
