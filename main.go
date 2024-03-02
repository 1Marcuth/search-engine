package main

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
	"strings"
	"unicode"
	"sort"
)

type Item struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type ItemSearched struct {
	*Item
	Relevance int `json:"relevance"`
}

type ItemsSearched []ItemSearched

type SearchEngine struct {
	items []Item
}

func (items ItemsSearched) Len() int {
	return len(items)
}

func (items ItemsSearched) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items ItemsSearched) Less(i, j int) bool {
	return items[i].Relevance > items[j].Relevance
}

func normalizeString(str string) string {
	transformed := trasnform("NFD", str)

	var stringBuilder strings.Builder

	for _, char := range transformed {
		if !unicode.Is(unicode.Mn, char) {
			stringBuilder.WriteRune(char)
		}
	}

	normalized := stringBuilder.String()

	var result strings.Builder
	
	for _, char := range normalized {
		if unicode.IsLetter(char) || unicode.IsNumber(char) || unicode.IsSpace(char) {
			result.WriteRune(char)
		}
	}

	return strings.ToLower(result.String())
}

func trasnform(from string, str string) string {
	transformed := make([]rune, 0, len(str))

	for _, char := range str {
		if char < 0 || char > unicode.MaxRune {
			continue
		}

		transformed = append(transformed, char)
	}

	return string(transformed)
}

func normalizeAndTokenizeString(str string) []string {
	normalizedString := normalizeString(str)
	tokens := strings.Fields(normalizedString)
	return tokens
}

func (searchEngine SearchEngine) Search(query string) []ItemSearched {
	var filteredItems []ItemSearched

	queryTokens := normalizeAndTokenizeString(query)

	for _, item := range searchEngine.items {
		itemSearched := ItemSearched{
			Item:      &item,
			Relevance: 0,
		}

		titleTokens := normalizeAndTokenizeString(item.Title)
		descriptionTokens := normalizeAndTokenizeString(item.Description)

		for _, queryToken := range queryTokens {
			for _, titleToken := range titleTokens {
				if titleToken == queryToken {
					itemSearched.Relevance += 30
				} else if strings.Contains(titleToken, queryToken) {
					itemSearched.Relevance += 25
				} else if strings.Contains(queryToken, titleToken) {
					itemSearched.Relevance += 20
				}
			}

			for _, descriptionToken := range descriptionTokens {
				if descriptionToken == queryToken {
					itemSearched.Relevance += 5
				} else if strings.Contains(descriptionToken, queryToken) {
					itemSearched.Relevance += 2
				} else if strings.Contains(queryToken, descriptionToken) {
					itemSearched.Relevance += 1
				}
			}

			for _, tag := range item.Tags {
				tagTokens := normalizeAndTokenizeString(tag)

				for _, tagToken := range tagTokens {
					if tagToken == queryToken {
						itemSearched.Relevance += 25
					} else if strings.Contains(tagToken, queryToken) {
						itemSearched.Relevance += 20
					} else if strings.Contains(queryToken, tagToken) {
						itemSearched.Relevance += 15
					}
				}
			}

			if itemSearched.Relevance > 0 {
				filteredItems = append(filteredItems, itemSearched)
			}
		}
	}

	return filteredItems
}

func main() {
	data, err := os.ReadFile("data.json")

	if err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v", err)
	}

	var items []Item

	if err := json.Unmarshal(data, &items); err != nil {
		log.Fatalf("Erro ao decodificao o arquivo JSON: %v", err)
	}

	searchEngine := SearchEngine{
		items: items,
	}

	resultItems := searchEngine.Search("floresta")
	sort.Sort(ItemsSearched(resultItems))

	jsonData, err := json.Marshal(resultItems)

	if err != nil {
		log.Fatalf("Erro ao codificar para JSON: %v", err)
	}

	err = os.WriteFile("output.json", jsonData, 0644)

	if err != nil {
		log.Fatalf("Erro ao escrever no arquivo: %v", err)
	}

	fmt.Println(string(jsonData))
}