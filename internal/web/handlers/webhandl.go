package webhandl

import (
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/shulganew/shear.git/internal/app/config"
	utils "github.com/shulganew/shear.git/internal/core"
	"github.com/shulganew/shear.git/internal/storage"
)

// hadler for  GET and POST  hor and log urls

type URLHandler struct {
	storage storage.URLSetGet
}

func (u *URLHandler) SetMapStorage(storage storage.URLSetGet) {
	u.storage = storage
}

func (u *URLHandler) GetStorage() storage.URLSetGet {
	return u.storage
}

// GET and redirect by shortUrl
func (u *URLHandler) GetURL(res http.ResponseWriter, req *http.Request) {

	shortURL := chi.URLParam(req, "id")

	//get long Url from storage
	longURL, exist := u.storage.GetLongURL(shortURL)

	//set content type
	res.Header().Add("Content-Type", "text/plain")
	log.Println("Redirect to: ", longURL)

	if exist {
		res.Header().Set("Location", longURL.String())
		//set status code 307
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}

}

// POTS and set generate short Url
func (u *URLHandler) SetURL(res http.ResponseWriter, req *http.Request) {

	readBody, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, "Body not found", http.StatusInternalServerError)
	}

	redirectURL, err := url.Parse(string(readBody))
	if err != nil {
		http.Error(res, "Wrong URL in body, parse error", http.StatusInternalServerError)
	}
	//from addres: from OS ENV
	config := config.GetConfig()

	//main URL = Shema + hostname + port
	mainURL := redirectURL.Scheme + "://" + config.ResultAddress

	shortURL := utils.GenerateShorLink()

	//join full long URL
	longStrURL, _ := url.JoinPath(mainURL, shortURL)
	longURL, _ := url.Parse(longStrURL)

	log.Println("Save long url: ", longURL)

	//save map to storage
	u.storage.SetURL(shortURL, *redirectURL)

	log.Println("Server ansver with short URL: ", longURL)

	//set content type
	res.Header().Add("Content-Type", "text/plain")

	//set status code 201
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte(longURL.String()))
}
