package main

import (
	"encoding/json"
	// "fmt"
	"io/ioutil"
	"log"
)

type Item struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Tags []string `json:"tags"`
}

type ItemSearched struct {
	*Item
	Relevance int
}

type SearchEngine struct {
	items []Item
}

func (self SearchEngine) Search(query string) []ItemSearched {
	var filteredItems []ItemSearched

	for _, item := range self.items {
		itemSearched := ItemSearched {
			Relevance: 0,
		}

		
	}

	return filteredItems
}

func main() {
	data, err := ioutil.ReadFile("data.json")

	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v", err)
	}

	var items []Item

	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatalf("Erro ao decodificao o arquivo JSON: %v", err)
	}

	searchEngine := SearchEngine {
		items: items,
	}

	searchEngine.Search()
}
