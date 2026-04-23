package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var cliCommands map[string]cliCommand
	cliCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Display a help message",
			callback:    func() error { return commandHelp(cliCommands) },
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
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
		if cmd, exists := cliCommands[command]; exists {
			if err := cmd.callback(); err != nil {
				fmt.Printf("Error executing command '%s': %v\n", command, err)
			}
		} else {
			fmt.Printf("Unknown command: '%s'\n", command)
		}
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n")
	for _, cmd := range commands {
		fmt.Printf("- %s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}
