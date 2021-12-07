package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Oauth2 struct {
	Status string `json:"status"`
	Result struct {
		Data struct {
			AccessToken string `json:"access_token"`
			ExpiresIn   int    `json:"expires_in"`
			TokenType   string `json:"token_type"`
			Scope       string `json:"scope"`
		} `json:"data"`
	} `json:"result"`
}

func AccessTokenFetch() string {
	url := "https://api.codechef.com/oauth/token"
	method := "POST"
	client_id := os.Getenv("CLIENT_ID")
	client_secret := os.Getenv("CLIENT_SECRET")
	redirect_uri := os.Getenv("REDIRECT_URI")

	payloadJson := map[string]string{
		"grant_type":    "client_credentials",
		"scope":         "public",
		"client_id":     client_id,
		"client_secret": client_secret,
		"redirect_uri":  redirect_uri}
	payload, err := json.Marshal(payloadJson)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	oauth2resp := Oauth2{}
	jsonErr := json.Unmarshal(body, &oauth2resp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return oauth2resp.Result.Data.AccessToken

}
func CodechefRankings(token string) interface{} {
	url := "https://api.codechef.com/ratings/all?fields=rating&country=India&institution=SRM%20Institute%20Of%20Science%20And%20Technology&sortBy=rating&sortOrder=desc&limit=100"
	method := "GET"

	client := &http.Client{
		Timeout: time.Second * 2,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var result interface{}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &result)
	return result
}

var Router *gin.Engine

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowCredentials: true,
	}))
	r.GET("/", func(c *gin.Context) {
		token := AccessTokenFetch()
		jsonRes := CodechefRankings(token)
		c.JSON(200, jsonRes)
	})
	r.Run()
}
