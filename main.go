package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	locations "mike_pok/internal"
	"os"
	"strings"
	"sync"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string, *Pokedex) error
}

type Config struct {
	Next     string
	Previous string
}

type Pokedex struct {
	mu      sync.Mutex
	Pokedex map[string]Pokemon
}

func NewPokedex() *Pokedex {
	return &Pokedex{
		Pokedex: make(map[string]Pokemon),
	}
}

func main() {

	config := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "",
	}

	registry := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapB,
		},
		"explore": {
			name:        "explore",
			description: "Lists all Pokemon in a single location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catching Pokemon adds them to the user's Pokedex.",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Prints the name, height, weight, stats and type(s) of a Pokemon.",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints a list of all pokemon in your pokedex",
			callback:    commandPokedex,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	pokedex := NewPokedex()

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		words := cleanInput(line)

		command := words[0]
		secondaryCommand := ""
		if len(words) > 1 {
			secondaryCommand = words[1]
		}

		value, ok := registry[command]

		if !ok {
			fmt.Println("Unknown command")
		} else {
			if err := value.callback(config, secondaryCommand, pokedex); err != nil {
				fmt.Println(err)
				if value.name == "help" {
					fmt.Println("usage:")
					for _, c := range registry {
						fmt.Printf("%v:%v \n", c.name, c.description)
					}
				}
				if value.name == "exit" {
					os.Exit(0)
				}

			}
		}

	}

}

func cleanInput(text string) []string {

	words := []string{}

	trimmedText := strings.TrimSpace(text)

	split_text := strings.Fields(trimmedText)

	for _, s := range split_text {
		words = append(words, (strings.ToLower(s)))
	}

	return words
}

func commandExit(config *Config, command string, pokedex *Pokedex) error {

	return errors.New("Closing the Pokedex... Goodbye!")
}

func commandHelp(config *Config, command string, pokedex *Pokedex) error {

	return errors.New("Welcome to the Pokedex!")
}

type Location struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Result   []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

func commandMap(config *Config, command string, pokedex *Pokedex) error {

	data := locations.GetPokemonResource(config.Next)

	location := Location{}

	err := json.Unmarshal(data, &location)

	if err != nil {
		log.Fatal(err)
	}

	config.Next = location.Next
	config.Previous = location.Previous

	for _, r := range location.Result {
		fmt.Println(r.Name)
	}

	return nil
}

func commandMapB(config *Config, command string, pokedex *Pokedex) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	data := locations.GetPokemonResource(config.Previous)

	location := Location{}

	err := json.Unmarshal(data, &location)

	if err != nil {
		log.Fatal(err)
	}

	config.Next = location.Next
	config.Previous = location.Previous

	for _, r := range location.Result {
		fmt.Println(r.Name)
	}

	return nil
}

type PokemonEncounter struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func commandExplore(config *Config, command string, pokedex *Pokedex) error {
	locationUrl := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%v", command)
	data := locations.GetPokemonResource(locationUrl)

	pokemonEncounters := PokemonEncounter{}

	err := json.Unmarshal(data, &pokemonEncounters)

	if err != nil {

		return errors.New("could not get pokemon encounters")
	}

	for _, pokemon := range pokemonEncounters.PokemonEncounters {
		fmt.Println(pokemon.Pokemon.Name)
	}

	return nil
}

type Pokemon struct {
	PokemonId      int    `json:"id"`
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func commandCatch(config *Config, command string, pokedex *Pokedex) error {
	fmt.Printf("Throwing a Pokeball at %v...", command)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%v", command)
	data := locations.GetPokemonResource(url)

	pokemon := Pokemon{}

	err := json.Unmarshal(data, &pokemon)

	if err != nil {
		return errors.New("could not get pokemon")
	}
	number := rand.Intn(pokemon.BaseExperience + 100)

	if number >= pokemon.BaseExperience {
		pokedex.mu.Lock()
		defer pokedex.mu.Unlock()

		pokedex.Pokedex[pokemon.Name] = pokemon
		fmt.Println(pokedex.Pokedex)
		fmt.Printf("%v was caught!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%v escaped!\n", pokemon.Name)

	return nil
}

func commandInspect(config *Config, command string, pokedex *Pokedex) error {
	fmt.Println(pokedex.Pokedex)
	pokemon, ok := pokedex.Pokedex[command]

	if !ok {
		fmt.Printf("You have not caught %v", command)
		return nil
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, typeInfo := range pokemon.Types {
		fmt.Println("  -", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedex(config *Config, command string, pokedex *Pokedex) error {

	fmt.Println("Your pokedex")
	for _, value := range pokedex.Pokedex {
		fmt.Printf("- %v\n", value.Name)
	}
	return nil
}
