package middlewares

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const PASS string = "mysecret"

func Cookie(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if userID, ok := getUserID(r); ok {
			//cookie iser_id is set
			cookies := r.Cookies()

			//clean cookie data
			r.Header["Cookie"] = make([]string, 0)
			for _, cookie := range cookies {

				if cookie.Name == "user_id" {
					cookie.Value = userID
				}
				r.AddCookie(cookie)
			}

			zap.S().Infoln("UserID: ", userID)

		} else {
			//cookie not set or not decoded
			//create new user uuid
			userID, err := uuid.NewV7()
			if err != nil {
				zap.S().Errorln("Error generate user uuid")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			//encode coockie for client
			coded, err := EncodeCookie(userID.String(), PASS)
			if err != nil {
				zap.S().Errorln("Error encode uuid")
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			//set to response
			codedCookie := http.Cookie{Name: "user_id", Value: coded}
			http.SetCookie(w, &codedCookie)

			//set to request
			cookie := http.Cookie{Name: "user_id", Value: userID.String()}
			r.AddCookie(&cookie)

		}

		h.ServeHTTP(w, r)

	})

}

func getUserID(r *http.Request) (userID string, ok bool) {
	//find UserID in cookies
	cookie, err := r.Cookie("user_id")
	if err != nil {
		return "", false
	}

	userID, err = DecodeCookie(cookie.Value, PASS)
	if err != nil {
		return "", false
	}

	//check correct UUID
	_, err = uuid.Parse(userID)
	if err != nil {
		return "", false
	}
	return userID, true

}

func DecodeCookie(secret string, password string) (uuid string, err error) {

	key := sha256.Sum256([]byte(password))
	msg, err := base64.StdEncoding.DecodeString(secret)

	if err != nil {
		return "", err
	}
	aesblock, err := aes.NewCipher(key[:32])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {

		return
	}

	lenth := aesgcm.NonceSize()
	nonceSize := len(key) - lenth
	nonce := key[nonceSize:]
	binary := []byte(msg)

	decrypted, err := aesgcm.Open(nil, nonce, binary, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	return string(decrypted), nil

}

func EncodeCookie(uuid string, password string) (secret string, err error) {
	key := sha256.Sum256([]byte(password))

	aesblock, err := aes.NewCipher(key[:32])
	if err != nil {

		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {

		return
	}

	lenth := aesgcm.NonceSize()
	nonceSize := len(key) - lenth
	nonce := key[nonceSize:]
	binary := []byte(uuid)

	coded := aesgcm.Seal(nil, nonce, binary, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	return base64.StdEncoding.EncodeToString(coded), nil

}
