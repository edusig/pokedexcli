package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const API_URL = "https://pokeapi.co/api/v2"

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func GetMapLocationAreas(urlParams string) (LocationAreaResponse, error) {
	locationResponse := LocationAreaResponse{}
	body := request("/location-area" + urlParams)
	err := json.Unmarshal(body, &locationResponse)
	if err != nil {
		return locationResponse, err
	}
	return locationResponse, nil
}

func request(path string) []byte {
	res, err := http.Get(API_URL + path)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	return body
}
