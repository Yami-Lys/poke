package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Yami-Lys/poke/internal/pokecache"
)

type config struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

func main() {
	cfg := &config{
		Cache:   pokecache.NewCache(5 * time.Minute),
		Pokedex: make(map[string]pokemon),
	}

	var cliCommands map[string]cliCommand
	cliCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    func(cfg *config, args []string) error { return commandHelp(cfg, cliCommands) },
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon in the current location area",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon in your Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all Pokemon in your Pokedex",
			callback:    commandPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		words := cleanInput(scanner.Text())

		if len(words) == 0 {
			continue
		}

		command := words[0]
		args := words[1:]
		if cmd, exists := cliCommands[command]; exists {
			if err := cmd.callback(cfg, args); err != nil {
				fmt.Printf("Error executing command '%s': %v\n", command, err)
			}
		} else {
			fmt.Printf("Unknown command: '%s'\n", command)
		}
	}
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, cmd := range commands {
		fmt.Printf("- %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}
