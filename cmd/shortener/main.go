package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
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
func getUrl(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {

		shortUrl := strings.TrimLeft(req.URL.String(), "/")
		longUrl, exist := Urldb[shortUrl]

		//set content type
		res.Header().Add("Content-Type", "text/plain")

		//set status code 307
		res.WriteHeader(http.StatusTemporaryRedirect)

		if exist {
			http.Redirect(res, req, longUrl, http.StatusTemporaryRedirect)
		}
		return

	} else if req.Method == http.MethodPost {

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
		fmt.Println(Urldb)

		res.Write([]byte(answer))
	}
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, getUrl)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}

//curl -v -H "Content-Type: text/plain" -X POST http://localhost:8080/ -d "{"https://yandex.ru"}"
//curl -v -H "Content-Type: text/plain" http://localhost:8080/hjnFtibr
