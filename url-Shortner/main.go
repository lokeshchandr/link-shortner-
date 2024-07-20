package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID          string    `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortURL    string    `json:"short_url"`
	Date        time.Time `json:"date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(OriginalURL))  //convet the url into the bites
	fmt.Println("hasher", hasher)
	data := hasher.Sum(nil)
	fmt.Println("hasher data0", data)

	hash := hex.EncodeToString(data)

	fmt.Println("encode to string ", hash)
	fmt.Println("final 8 bit string is", hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		Date:        time.Now(),
	}
	return shortURL
}
func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL NOT FOUND")
	}

	return url, nil
}
func RootPageURL(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func ShortUrlHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid REquiest body", http.StatusBadRequest)
		return
	}
	shortURL_ := createURL(data.URL)

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL_}

	w.Header().Set("content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]

	url, err := getURL(id)
	if err != nil {
		http.Error(w, "invalide responce ", http.StatusNotFound)
	
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// fmt.Println("urlshoetner")
	// originalURL := "localhost 125.0.0.1"
	// generateShortURL(originalURL)

	http.HandleFunc("/", RootPageURL)
	http.HandleFunc("/shorten", ShortUrlHandler)
	http.HandleFunc("/redirect/",redirectURLHandler)

	fmt.Println("server started")
	err := http.ListenAndServe(":3000", nil) //to start the server
	if err != nil {
		fmt.Println("Erorr found", err)
	}

}
