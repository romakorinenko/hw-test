package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountWords(t *testing.T) {
	t.Run("В метод передан комбинированный текст, включающий все предыдущие кейсы", func(t *testing.T) {
		text := "кот!@#$%^&*()_+=-?><}{[]|\"'Кот кОт кот1рыжий:рыжий*8Hello,世界-🤯🤯🤯 ta4ка"

		expectedResult := map[string]int{"кот": 3, "кот1рыжий": 1, "рыжий": 1, "8hello": 1, "世界": 1, "🤯🤯🤯": 1, "ta4ка": 1}
		words := countWords(text)

		require.Equal(t, expectedResult, words)
	})
}
