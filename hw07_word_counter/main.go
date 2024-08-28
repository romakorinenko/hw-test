package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(countWords("Я устал, босс"))
}

func countWords(text string) map[string]int {
	textWithoutPunctuationRegex := regexp.MustCompile(`[[:punct:]]`)
	textWithoutPunctuation := textWithoutPunctuationRegex.ReplaceAllString(text, " ")

	wordRegex := regexp.MustCompile(`\S+`)
	words := wordRegex.FindAllString(strings.ToLower(textWithoutPunctuation), -1)

	wordToCountMap := make(map[string]int)

	for _, word := range words {
		wordToCountMap[word]++
	}

	return wordToCountMap
}
