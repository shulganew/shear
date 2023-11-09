package webhandl

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/shulganew/shear.git/internal/app/config"
	utils "github.com/shulganew/shear.git/internal/core"
	"github.com/shulganew/shear.git/internal/storage"
)

// hadler for  GET and POST  hor and log urls

// GET and redirect by shorUrl
func GetUrl(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("GET")
	shortUrl := strings.TrimLeft(req.URL.String(), "/")

	urldb := storage.GetUrldb()
	longUrl, exist := (*urldb)[shortUrl]

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 307
	//res.WriteHeader(http.StatusTemporaryRedirect)

	//set status code 307
	log.Println("Redirect to: ", longUrl)

	if exist {
		//res.Header().Add("Location", "example.com")
		http.Redirect(res, req, longUrl, http.StatusTemporaryRedirect)
	}
}

// POTS and set generate short Url
func SetUrl(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("POTS")
	//answer := fmt.Sprintf("Method: %s\r\n", req.Method)
	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	//from addres: from OS ENV
	config := config.GetConfig()
	answer := "http://" + config.ResultAddress + "/"

	longUrl := string(readBody)
	log.Println("Save long url: ", longUrl)
	shortUrl := utils.GenerateShorLink()

	//save map to storage
	urldb := storage.GetUrldb()

	(*urldb)[shortUrl] = longUrl

	answer += shortUrl

	log.Println("Server ansver with short URL: ", answer)
	//log.Println("Config POTS: ", configShear)
	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(answer))
}
