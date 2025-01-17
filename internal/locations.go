package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func GetLocations(url string) []byte {
	cache := &Cache{}

	existingCache, cacheExists := cache.Get(url)

	if !cacheExists {
		fmt.Println("Cache does not exist for this url")
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

		cache := NewCache(5 * time.Second)
		cache.Add(url, body)

		return body
	}

	return existingCache
}
