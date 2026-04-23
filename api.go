package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2"

type locationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func fetchLocationAreas(url string, cfg *config) (locationAreaResponse, error) {
	if val, ok := cfg.Cache.Get(url); ok {
		fmt.Println("(cache hit)")
		var locationResp locationAreaResponse
		if err := json.Unmarshal(val, &locationResp); err != nil {
			return locationAreaResponse{}, fmt.Errorf("failed to parse cached response: %w", err)
		}
		return locationResp, nil
	}

	fmt.Println("(fetching from API...)")
	resp, err := http.Get(url)
	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("failed to fetch location areas: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return locationAreaResponse{}, fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.Cache.Add(url, body)

	var locationResp locationAreaResponse
	if err := json.Unmarshal(body, &locationResp); err != nil {
		return locationAreaResponse{}, fmt.Errorf("failed to parse response: %w", err)
	}

	return locationResp, nil
}

func commandMap(cfg *config) error {
	url := baseURL + "/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	locationResp, err := fetchLocationAreas(url, cfg)
	if err != nil {
		return err
	}

	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	locationResp, err := fetchLocationAreas(*cfg.Previous, cfg)
	if err != nil {
		return err
	}

	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}
