package main

import (
	"bufio"
	"errors"
	"fmt"
	"internal/pokeapi"
	"internal/pokecache"
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

type cliCommandConfig struct {
	next            string
	prev            string
	apiClient       *pokeapi.APIClient
	capturedPokemon map[string]Pokemon
}

type Pokemon struct {
	name           string
	id             int
	baseExperience int
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cliCommandConfig, ...string) error
}

func getPointerValueOrDefault[T any](val *T, def T) T {
	if val == nil {
		return def
	}
	return *val
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
		"explore": {
			name:        "explore <area_name>",
			description: "Explore an area and find the pokemon that resides there.",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Tries to catch a pokemon",
			callback:    commandCatchPokemon,
		},
	}
}

func commandHelp(*cliCommandConfig, ...string) error {
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

func commandExit(*cliCommandConfig, ...string) error {
	os.Exit(0)
	return nil
}

func commandMap(config *cliCommandConfig, _ ...string) error {
	urlParams := ""
	if config.next != "" {
		parsedNext, err := url.Parse(config.next)
		if err == nil {
			urlParams = "?" + parsedNext.RawQuery
		}
	}
	locationAreas, err := config.apiClient.GetMapLocationAreas(urlParams)
	if err != nil {
		fmt.Println("Could not get the locations from the API. Try again later.")
	}

	config.next = getPointerValueOrDefault(locationAreas.Next, "")
	config.prev = getPointerValueOrDefault(locationAreas.Previous, "")

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapB(config *cliCommandConfig, _ ...string) error {
	urlParams := ""
	if config.prev != "" {
		parsedPrev, err := url.Parse(config.prev)
		if err == nil {
			urlParams = "?" + parsedPrev.RawQuery
		}
	}
	locationAreas, err := config.apiClient.GetMapLocationAreas(urlParams)
	if err != nil {
		return errors.New("could not get the locations from the API. Try again later")
	}

	config.next = getPointerValueOrDefault(locationAreas.Next, "")
	config.prev = getPointerValueOrDefault(locationAreas.Previous, "")

	for _, location := range locationAreas.Results {
		fmt.Println(location.Name)
	}
	return nil
}

func commandExplore(config *cliCommandConfig, args ...string) error {
	if len(args) < 1 {
		return errors.New("missing location area name")
	}
	area := args[0]
	locationArea, err := config.apiClient.GetMapLocationArea(area)
	if err != nil {
		return errors.New("could not get location detail. Try again later")
	}

	fmt.Printf("Exploring %v\n", area)
	fmt.Printf("Found %v Pokemon:\n", len(locationArea.PokemonEncounters))
	for _, pokemon := range locationArea.PokemonEncounters {
		fmt.Printf("- %v\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatchPokemon(config *cliCommandConfig, args ...string) error {
	if len(args) < 1 {
		return errors.New("missing pokemon name")
	}
	pokemonName := args[0]
	pokemon, err := config.apiClient.GetPokemonDetail(pokemonName)
	if err != nil {
		return errors.New("could not get pokemon detail. Try again later")
	}

	fmt.Printf("Throwing a Pokeball at %v...\n", pokemonName)

	maxDifficulty := 900
	minDifficulty := 190
	randDifficulty := rand.Intn(maxDifficulty - pokemon.BaseExperience - minDifficulty)
	totalDifficulty := minDifficulty + randDifficulty + pokemon.BaseExperience

	maxChance := 1000
	randChance := rand.Intn(maxChance)

	if randChance > totalDifficulty {
		fmt.Printf("%v was caught!\n", pokemon.Name)
		config.capturedPokemon[pokemon.Name] = Pokemon{
			name:           pokemon.Name,
			id:             pokemon.ID,
			baseExperience: pokemon.BaseExperience,
		}
	} else {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cmds := getCliCommands()
	cache := pokecache.NewCache(30 * time.Second)
	client := pokeapi.NewAPIClient(cache)
	config := cliCommandConfig{apiClient: &client, capturedPokemon: make(map[string]Pokemon)}
	for {
		fmt.Print("pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		words := strings.Split(line, " ")
		command := words[0]
		if cmd, ok := cmds[command]; ok {
			err := cmd.callback(&config, words[1:]...)
			if err != nil {
				fmt.Printf("Command error: %v", err)
				commandHelp(&config)
			}
		} else {
			commandHelp(&config)
		}
	}
}
