package main

import (
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

var Urldb = map[string]string{}

//generate short link

func generateShorLink() string {

	//base charset
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	//nuber of short chars in url string
	n := 8

	sb := strings.Builder{}
	sb.Grow(7)
	for i := 0; i < n; i++ {
		sb.WriteByte(charset[rand.Intn(len(charset))])
	}
	return sb.String()
}

// hadler for  GET and POST  hor and log urls

// GET and redirect by shorUrl
func getUrl(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("GET")
	shortUrl := strings.TrimLeft(req.URL.String(), "/")
	longUrl, exist := Urldb[shortUrl]

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 307
	//res.WriteHeader(http.StatusTemporaryRedirect)

	//set status code 307
	if exist {
		http.Redirect(res, req, longUrl, http.StatusTemporaryRedirect)
	}
}

// POTS and set generate short Url
func setUrl(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("POTS")
	//answer := fmt.Sprintf("Method: %s\r\n", req.Method)
	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	//from addres
	answer := "http://" + req.Host + "/"

	longUrl := string(readBody)
	shortUrl := generateShorLink()

	Urldb[shortUrl] = longUrl
	answer += shortUrl

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)

	//remove after tests
	//fmt.Println(Urldb)

	res.Write([]byte(answer))
}

func main() {
	r := chi.NewRouter()
	r.Get("/{id}", getUrl)
	r.Post("/", setUrl)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}

//curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/ -d "{"https://yandex.ru"}"
//curl -v -H "Content-Type: text/plain" http://localhost:8080/hjnFtibr
