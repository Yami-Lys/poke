package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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

type exploreResponse struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Type []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon-name>")
	}
	name := args[0]
	url := baseURL + "/pokemon/" + name
	fmt.Printf("Throwing a Pokeball at %s...\n", name)

	body, err := fetchURL(url, cfg)
	if err != nil {
		return fmt.Errorf("pokemon not found: %s", name)
	}

	var p pokemon
	if err := json.Unmarshal(body, &p); err != nil {
		return fmt.Errorf("failed to parse pokemon data: %w", err)
	}

	catchChance := 1.0 - (float64(p.BaseExperience) / float64(p.BaseExperience+100))
	if rand.Float64() < catchChance {
		fmt.Printf("%s was caught!\n", name)
		cfg.Pokedex[name] = p
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func fetchURL(url string, cfg *config) ([]byte, error) {
	if val, ok := cfg.Cache.Get(url); ok {
		return val, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	cfg.Cache.Add(url, body)
	return body, nil
}

func commandMap(cfg *config, args []string) error {
	url := baseURL + "/location-area"
	if cfg.Next != nil {
		url = *cfg.Next
	}

	body, err := fetchURL(url, cfg)
	if err != nil {
		return err
	}

	var locationResp locationAreaResponse
	if err := json.Unmarshal(body, &locationResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	body, err := fetchURL(*cfg.Previous, cfg)
	if err != nil {
		return err
	}

	var locationResp locationAreaResponse
	if err := json.Unmarshal(body, &locationResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	cfg.Next = locationResp.Next
	cfg.Previous = locationResp.Previous

	for _, area := range locationResp.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <area-name>")
	}

	area := args[0]
	url := baseURL + "/location-area/" + area
	fmt.Printf("Exploring %s...\n", area)

	body, err := fetchURL(url, cfg)
	if err != nil {
		return err
	}

	var exploreResp exploreResponse
	if err := json.Unmarshal(body, &exploreResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range exploreResp.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon-name>")
	}
	name := args[0]
	p, exists := cfg.Pokedex[name]
	if !exists {
		return fmt.Errorf("%s is not in your Pokedex. Try catching it first!", name)
	}

	fmt.Printf("Name: %s\n", p.Name)
	fmt.Printf("Base Experience: %d\n", p.BaseExperience)
	fmt.Printf("Height: %d\n", p.Height)
	fmt.Printf("Weight: %d\n", p.Weight)
	fmt.Println("Stats:")
	for _, stat := range p.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range p.Type {
		fmt.Printf(" - %s\n", t.Type.Name)
	}
	return nil
}
