package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	locations "mike_pok/internal"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
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
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		words := cleanInput(line)

		command := words[0]

		value, ok := registry[command]

		if !ok {
			fmt.Println("Unknown command")
		} else {
			if err := value.callback(config); err != nil {
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

func commandExit(config *Config) error {

	return errors.New("Closing the Pokedex... Goodbye!")
}

func commandHelp(config *Config) error {

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

func commandMap(config *Config) error {

	data := locations.GetLocations(config.Next)

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

func commandMapB(config *Config) error {
	if config.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	data := locations.GetLocations(config.Previous)

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
