package main

import (
	"fmt"
	"regexp"
	"strings"
)

func main() {
	fmt.Println(countWords("кот Кот кОт кот1рыжий:рыжий*8Hello,世界"))
}

func countWords(text string) map[string]int {
	deleteAllInsteadLetters := regexp.MustCompile(`[^\p{L}]+`)
	onlyWordsString := deleteAllInsteadLetters.ReplaceAllString(text, " ")

	wordRegex := regexp.MustCompile(`\S+`)
	words := wordRegex.FindAllString(strings.ToLower(onlyWordsString), -1)

	wordToCountMap := make(map[string]int)

	for _, word := range words {
		wordInLowerCase := strings.ToLower(word)
		wordToCountMap[wordInLowerCase]++
	}

	return wordToCountMap
}
