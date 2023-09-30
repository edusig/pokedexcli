package main

import (
	"bufio"
	"fmt"
	"internal/pokeapi"
	"net/url"
	"os"
)

type cliCommandConfig struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cliCommandConfig) error
}

func getCliCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays names of 20 location areas. Each subsequent call displays the next 20",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Similar to map, however displays the previous 20 location areas",
			callback:    commandMapB,
		},
	}
}

func commandHelp(*cliCommandConfig) error {
	cmds := getCliCommands()
	fmt.Printf(`
Welcome to the Pokedex!
Usage:

`)
	for _, val := range cmds {
		fmt.Printf("%v: %v\n", val.name, val.description)
	}
	fmt.Println()
	return nil
}

func commandExit(*cliCommandConfig) error {
	os.Exit(0)
	return nil
}

func commandMap(config *cliCommandConfig) error {
	urlParams := ""
	if config.Next != "" {
		parsedNext, err := url.Parse(config.Next)
		if err == nil {
			urlParams = "?" + parsedNext.RawQuery
		}
	}
	locationAreas, err := pokeapi.GetMapLocationAreas(urlParams)
	if err != nil {
		fmt.Println("Could not get the locations from the API. Try again later.")
	}

	if locationAreas.Next != nil {
		config.Next = *locationAreas.Next
	}
	if locationAreas.Previous != nil {
		config.Previous = *locationAreas.Previous
	}

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(config *cliCommandConfig) error {
	urlParams := ""
	if config.Previous != "" {
		parsedPrev, err := url.Parse(config.Previous)
		if err == nil {
			urlParams = "?" + parsedPrev.RawQuery
		}
	}
	locationAreas, err := pokeapi.GetMapLocationAreas(urlParams)
	if err != nil {
		fmt.Println("Could not get the locations from the API. Try again later.")
	}

	if locationAreas.Next != nil {
		config.Next = *locationAreas.Next
	}
	if locationAreas.Previous != nil {
		config.Previous = *locationAreas.Previous
	}

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := getCliCommands()
	config := cliCommandConfig{}
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		if cmd, ok := cmds[line]; ok {
			err := cmd.callback(&config)
			if err != nil {
				fmt.Println("Error while executing command")
				break
			}
		} else {
			commandHelp(&config)
		}
	}
}
