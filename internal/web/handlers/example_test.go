package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Example show how to get short (brief) URL from shortener. Use POST request with original URL in the body.
func ExampleHandlerURL_SetURL() {
	// Create a Resty Client
	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "plain/text").
		// Add full URL to body
		SetBody("http://yandex.ru").
		Post("localhost:8080/")
	// Check error.
	if err != nil {
		panic(err)
	}
	// Get short URL.
	// Explore response object.
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body, short URL:\n", resp)
	fmt.Println()
}

// Example show how to get original URL from shortener. Use GET request with original URL in the body.
func ExampleHandlerURL_GetURL() {
	// Create a Resty Client with redirect checking.
	client := resty.New().SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		// Shortener make redirect response to original URL, check redirect in the response.
		for _, resp := range via {
			fmt.Println("Next host URL:", resp.URL)
		}
		return nil
	}))
	_, err := client.R().
		SetHeader("Content-Type", "plain/text").
		Get("http://localhost:8080/kMerdbZY")

	// Check error.
	if err != nil {
		panic(err)
	}

}

// Create several shorts with batch request API.
func ExampleHandlerBatch_BatchSet() {
	client := resty.New()
	resp, err := client.R().
		SetBody(`
		[
			{
			    "correlation_id": "id1",
			    "original_url": "http://yandex11.ru"
			},
			{
			    "correlation_id": "id2",
			    "original_url": "http://yandex12.ru"
			},
			{
			    "correlation_id": "id3",
			    "original_url": "http://yandex13.ru"
			}
		  ]
		`).
		SetHeader("Content-Type", "application/json").
		Post("http://localhost:8080/api/shorten/batch")

	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
}

// API for deleting user's URl. Send array in the body with briefs (short URLs) and add user id to the cookie.
func ExampleDelShorts_DelUserURLs() {
	client := resty.New()
	// Add cookie - user id for particular users URL delete. Cookie created by server during new user's URLs creation.
	client.SetCookie(&http.Cookie{
		Name:  "user_id",
		Value: "lsjnN107/XV2f27GIedhnl3eRKQPwspek3FoFK3AcYkRPxlhY9hN3DGxG9jQSBlZLo51mg==",
	})
	resp, err := client.R().
		SetBody(`
		["HNWPvptC", "AoloIhve", "jNVNooSF"]
		`).
		SetHeader("Content-Type", "plain/text").
		Delete("http://localhost:8080/api/user/urls")

	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()
}

// JSON API for getting short URL by original URL.
func ExampleHandlerAPI_GetBrief() {
	client := resty.New()
	resp, err := client.R().
		SetBody(`
		{
			"url": "https://practicum.yandex.ru"
		} 
		`).
		SetHeader("Content-Type", "application/json").
		Post("http://localhost:8080/api/shorten")

	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	//Response example:
	//
	// {"result":"https://localhost:8080/XsVEJmMy"}
}

// Get all user's original and short URLs API.
func ExampleHandlerAuth_GetUserURLs() {
	// Create a Resty Client with redirect checking.
	client := resty.New()
	// Add cookie - user id for particular users URLs. Cookie created by server during new user's URLs creation.
	client.SetCookie(&http.Cookie{
		Name:  "user_id",
		Value: "lsjnN107/XV2f27GIedhnl3eRKQPwspek3FoFK3AcYkRPxlhY9hN3DGxG9jQSBlZLo51mg==",
	})
	resp, err := client.R().
		Get("http://localhost:8080/api/user/urls")

	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Proto      :", resp.Proto())
	fmt.Println("  Time       :", resp.Time())
	fmt.Println("  Received At:", resp.ReceivedAt())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Response example:
	//
	//		[
	//	    {
	//	        "short_url": "http://localhost:8080/hyGZvXio",
	//	        "original_url": "http://yandex24.ru"
	//	    },
	//	    {
	//	        "short_url": "http://localhost:8080/jNVNooSF",
	//	        "original_url": "http://yandex31.ru"
	//	    },
	//	    {
	//	        "short_url": "http://localhost:8080/AoloIhve",
	//	        "original_url": "http://yandex32.ru"
	//	    },
	//	    {
	//	        "short_url": "http://localhost:8080/HNWPvptC",
	//	        "original_url": "http://yandex33.ru"
	//	    }
	//
	// ]
}
