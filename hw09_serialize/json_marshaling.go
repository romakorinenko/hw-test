package main

import (
	"encoding/json"
)

func (b *Book) MarshalJSON() ([]byte, error) {
	type Alias Book
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}

func (b *Book) UnmarshalJSON(bytes []byte) error {
	type Alias Book
	aliasStruct := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	return json.Unmarshal(bytes, aliasStruct)
}

func MarshalBooksJSON(books []*Book) ([]byte, error) {
	return json.Marshal(books)
}

func UnmarshalBooksJSON(data []byte) ([]*Book, error) {
	var books []*Book
	err := json.Unmarshal(data, &books)
	return books, err
}
