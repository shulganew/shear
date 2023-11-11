package webhandl

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/shulganew/shear.git/internal/app/config"
	utils "github.com/shulganew/shear.git/internal/core"
	"github.com/shulganew/shear.git/internal/storage"
)

// hadler for  GET and POST  hor and log urls

// GET and redirect by shorUrl
func GetUrl(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("GET")
	shortUrl := chi.URLParam(req, "id")
	log.Println("short URL: ", shortUrl)

	urldb := storage.GetUrldb()
	//get long Url from storage
	longUrl, exist := (*urldb)[shortUrl]

	//set content type
	res.Header().Add("Content-Type", "text/plain")
	log.Println("Redirect to: ", longUrl)

	if exist {
		res.Header().Set("Location", longUrl.String())
	} else {
		res.WriteHeader(http.StatusNotFound)
	}

	//set status code 307
	res.WriteHeader(http.StatusTemporaryRedirect)
}

// POTS and set generate short Url
func SetUrl(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("POTS")
	//answer := fmt.Sprintf("Method: %s\r\n", req.Method)
	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	redirectUrl, err := url.Parse(string(readBody))
	if err != nil {
		http.Error(res, "Wrong URL in body, parse error", http.StatusInternalServerError)
	}
	//from addres: from OS ENV
	config := config.GetConfig()
	//main URL = Shema + hostname + port
	mainUrl := redirectUrl.Scheme + "://" + config.ResultAddress

	shortUrl := utils.GenerateShorLink()

	//join full long URL
	longStrUrl, _ := url.JoinPath(mainUrl, shortUrl)
	longUrl, _ := url.Parse(longStrUrl)

	log.Println("Save long url: ", longUrl)

	//save map to storage
	urldb := storage.GetUrldb()
	(*urldb)[shortUrl] = *redirectUrl

	log.Println("Server ansver with short URL: ", longUrl)
	//log.Println("Config POTS: ", configShear)
	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(longUrl.String()))
}
