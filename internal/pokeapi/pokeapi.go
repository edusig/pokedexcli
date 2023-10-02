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

type Cache interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

type APIClient struct {
	cache Cache
}

func (api *APIClient) GetMapLocationAreas(urlParams string) (LocationAreaResponse, error) {
	locationResponse := LocationAreaResponse{}
	body := api.request("/location-area" + urlParams)
	err := json.Unmarshal(body, &locationResponse)
	if err != nil {
		return locationResponse, err
	}
	return locationResponse, nil
}

func (api *APIClient) request(path string) []byte {
	url := API_URL + path
	if api.cache != nil {
		if val, ok := api.cache.Get(url); ok {
			return val
		}
	}
	res, err := http.Get(url)
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
	api.cache.Add(url, body)
	return body
}

func NewAPIClient(cache Cache) APIClient {
	return APIClient{cache: cache}
}
