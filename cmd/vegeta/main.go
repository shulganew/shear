package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	vegeta "github.com/tsenart/vegeta/lib"
)

const userID = "lsjnN107/XV2f27GIedhnl3eRKQPwspek3FoFK3AcYkRPxlhY9hN3DGxG9jQSBlZLo51mg=="
const testDuration = 30 * time.Second
const numOfURLs = 5

// Load handler.
//
//	/ [post]
//	/{id} [get]
func SetURL() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		// Adding with resty. Getting using vegeta.
		client := resty.New()
		client.SetCookie(&http.Cookie{
			Name:  "user_id",
			Value: userID,
		})

		rand := rand.Intn(200000)
		payload := "http://ya" + strconv.Itoa(rand) + ".localhost"

		resp, err := client.R().
			SetHeader("Content-Type", "plain/text").

			// Add full URL to body.
			SetBody(payload).
			Post("http://localhost:8080")
		// Check error.
		if err != nil {
			panic(err)
		}

		// GET URL.
		if tgt == nil {
			return vegeta.ErrNilTarget
		}
		tgt.Method = http.MethodGet
		tgt.URL = string(resp.Body())
		header := http.Header{}
		header.Set("Cookie", "user_id="+userID)
		header.Set("Content-type", "plain/text")
		tgt.Header = header
		tgt.Body = []byte(payload)
		return nil
	}
}

// Load handler.
//
//	/api/user/urls [get]
func UsersUrls() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		for i := 0; i < numOfURLs; i++ {
			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "user_id",
				Value: userID,
			})
			rand := rand.Intn(2000000)
			payload := "http://ya" + strconv.Itoa(rand) + ".localhost"

			_, err := client.R().
				SetHeader("Content-Type", "plain/text").
				// Add full URL to body.
				SetBody(payload).
				Post("http://localhost:8080")
			// Check error.
			if err != nil {
				panic(err)
			}

		}

		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		tgt.Method = http.MethodGet
		tgt.URL = "http://localhost:8080/api/user/urls"
		header := http.Header{}
		header.Set("Cookie", "user_id="+userID)
		tgt.Header = header
		return nil
	}
}

// Load handler.
//
//	/api/shorten/batch [delete]
func DelUsersUrls() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		// Array of short URLs for Users delete.
		shorts := make([]string, 0)
		for i := 0; i < numOfURLs; i++ {
			client := resty.New()
			client.SetCookie(&http.Cookie{
				Name:  "user_id",
				Value: userID,
			})
			rand := rand.Intn(200000)
			payload := "http://ya.del" + strconv.Itoa(rand) + ".localhost"

			resp, err := client.R().
				SetHeader("Content-Type", "plain/text").
				// Add full URL to body.
				SetBody(payload).
				Post("http://localhost:8080")
			// Check error.
			if err != nil {
				panic(err)
			}
			// save ansver to array of short URL
			rawAnswerURL := string(resp.Body())
			answer := strings.Trim(rawAnswerURL, "\n")
			URL, err := url.Parse(answer)
			if err != nil {
				fmt.Println("Error: ", string(resp.Body()))
				panic(err)
			}
			path := URL.Path
			short := strings.TrimPrefix(path, "/")
			shorts = append(shorts, short)
		}

		if tgt == nil {
			return vegeta.ErrNilTarget
		}

		body, err := json.Marshal(&shorts)
		if err != nil {
			panic(err)
		}
		tgt.Method = http.MethodDelete
		tgt.URL = "http://localhost:8080/api/user/urls"
		header := http.Header{}
		header.Set("Cookie", "user_id="+userID)
		header.Set("Content-type", "plain/text")
		tgt.Header = header
		tgt.Body = []byte(body)
		return nil
	}
}

// Load handler.
//
//	/api/shorten [post]
func JSON() vegeta.Targeter {
	return func(tgt *vegeta.Target) error {
		// Array of short URLs for Users delete.
		tgt.Method = http.MethodPost
		tgt.URL = "http://localhost:8080/api/shorten"
		header := http.Header{}
		header.Set("Cookie", "user_id="+userID)
		header.Set("Content-type", "application/json")
		tgt.Header = header
		rand := rand.Intn(200000)
		payload := "{ \"url\": \"http://ya" + strconv.Itoa(rand) + "\" }"
		tgt.Body = []byte(payload)
		return nil
	}
}

func main() {
	targeters := make([]vegeta.Targeter, 0)
	targeters = append(targeters, SetURL())
	targeters = append(targeters, UsersUrls())
	targeters = append(targeters, DelUsersUrls())
	targeters = append(targeters, JSON())

	var s sync.WaitGroup
	for _, target := range targeters {
		s.Add(1)
		atak(&s, target)
	}
	s.Wait()
	fmt.Println("Done!")
}

func atak(s *sync.WaitGroup, tr vegeta.Targeter) {
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics
	rate := vegeta.Rate{Freq: 100, Per: 2 * time.Second}
	for res := range attacker.Attack(tr, rate, testDuration, reflect.TypeOf(tr).String()) {
		metrics.Add(res)
	}
	metrics.Close()
	fmt.Println("Latencies: ", metrics.Latencies)
	fmt.Println("Earliest: ", metrics.Earliest)
	fmt.Println("BytesIn: ", metrics.BytesIn)
	fmt.Println("BytesOut: ", metrics.BytesOut)
	fmt.Println("Rate: ", metrics.Rate)
	fmt.Println("Errors: ", metrics.Errors)
	fmt.Println("Success: ", metrics.Success)
	s.Done()
}
