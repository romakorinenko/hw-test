package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(countWords("кот Кот кОт кот1рыжий:рыжий"))
}

func countWords(text string) map[string]int {
	regex := regexp.MustCompile(`[a-zA-Zа-яА-я]+`)
	words := regex.FindAllString(text, -1)

	wordToCountMap := make(map[string]int)

	for _, word := range words {
		wordInLowerCase := strings.ToLower(word)
		count, ok := wordToCountMap[wordInLowerCase]
		if ok {
			wordToCountMap[wordInLowerCase] = count + 1
		} else {
			wordToCountMap[wordInLowerCase] = 1
		}
	}

	return wordToCountMap
}
