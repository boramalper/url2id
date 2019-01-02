package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"time"
)

type result struct {
	Doi string `json:"doi,omitempty"`
	Pmid string `json:"pmid,omitempty"`
	Pmcid string `json:"pmcid,omitempty"`
	Nihmsid string `json:"nihmsid,omitempty"`
	Emsid string `json:"emsid,omitempty"`
	Camsid string `json:"camsid,omitempty"`
}

// Regular Expressions to find identifiers in the webpage.
var doiRE, pmcidRE, nihmsidRE, emsidRE, camsidRE *regexp.Regexp

func main() {
	// Initialise regexes
	doiRE = regexp.MustCompile(`doi\.org/([^"]+)"`)
	pmcidRE = regexp.MustCompile(`(?i)PMC(\d+)`)
	nihmsidRE = regexp.MustCompile(`(?i)NIHMS(\d+)`)
	emsidRE = regexp.MustCompile(`(?i)EMS(\d+)`)
	camsidRE = regexp.MustCompile(`(?i)cams(\d+)`)

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
	var result result

	urls, ok := r.URL.Query()["url"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("url parameter is missing"))
		return
	}

	url_, err := url.Parse(urls[0])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("url cannot be parsed"))
		return
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("error with initialising HTTP client"))
		log.Printf("ERROR: cookiejar.New: %s", err.Error())
		return
	}

	client := &http.Client{
		Jar: jar,
		Timeout: 10 * time.Second,
	}

	req, _ := http.NewRequest("GET", urls[0], nil)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Host", url_.Host)
	req.Header.Set("User-Agent", "what")

	response, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not fetch the resource"))
		log.Printf("ERROR: http.Get(%s): %s", urls[0], err.Error())
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not read the body"))
		log.Printf("ERROR: readAll(%s): %s", urls[0], err.Error())
		return
	}

	if matches := doiRE.FindSubmatch(body); matches != nil {
		result.Doi, err = url.QueryUnescape(string(matches[1]))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("could not decode DOI"))
			log.Printf("ERROR: QueryUnescape(%s): %s", matches[1], urls[0])
			return
		}
	}

	if matches := pmcidRE.FindSubmatch(body); matches != nil {
		result.Pmcid = "PMC" + string(matches[1])
	}

	if matches := nihmsidRE.FindSubmatch(body); matches != nil {
		result.Nihmsid = "NIHMS" + string(matches[1])
	}

	if matches := emsidRE.FindSubmatch(body); matches != nil {
		result.Emsid = "EMS" + string(matches[1])
	}

	if matches := camsidRE.FindSubmatch(body); matches != nil {
		result.Camsid = "cams" + string(matches[1])
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=21600")
	_ = json.NewEncoder(w).Encode(result) // ignore errors
}
