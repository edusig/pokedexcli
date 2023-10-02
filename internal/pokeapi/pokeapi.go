package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const API_URL = "https://pokeapi.co/api/v2"

type Cache interface {
	Add(key string, val []byte)
	Get(key string) ([]byte, bool)
}

type APIClient struct {
	cache Cache
}

func NewAPIClient(cache Cache) APIClient {
	return APIClient{cache: cache}
}

func (api *APIClient) GetMapLocationAreas(urlParams string) (LocationAreasResponse, error) {
	resp := LocationAreasResponse{}
	body := api.request("/location-area" + urlParams)
	return unmarshalResponse(resp, body)
}

type LocationAreasResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (api *APIClient) GetMapLocationArea(area string) (LocationAreaResponse, error) {
	resp := LocationAreaResponse{}
	body := api.request("/location-area/" + area)
	return unmarshalResponse(resp, body)
}

type LocationAreaResponse struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (api *APIClient) GetPokemonDetail(pokemonName string) (PokemonDetailResponse, error) {
	resp := PokemonDetailResponse{}
	body := api.request("/pokemon/" + pokemonName)
	return unmarshalResponse(resp, body)
}

type PokemonDetailResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

func unmarshalResponse[T any](resp T, body []byte) (T, error) {
	err := json.Unmarshal(body, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
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
