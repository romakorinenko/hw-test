package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBook_MarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		book         *Book
		expectedJSON string
	}{
		{
			name: "все поля заполнены",
			book: &Book{
				Id:     1,
				Title:  "book",
				Author: "author",
				Year:   2029,
				Size:   300,
				Rate:   4.33,
			},
			expectedJSON: `{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33}`,
		},
		{
			name:         "все поля не заполнены",
			book:         &Book{},
			expectedJSON: "{}",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bytes, err := testCase.book.MarshalJSON()
			require.NoError(t, err)
			require.Equal(t, testCase.expectedJSON, string(bytes))
		})
	}
}

func TestBook_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name         string
		bookBytes    []byte
		expectedBook *Book
	}{
		{
			name:      "json содержит все поля книги",
			bookBytes: []byte(`{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33}`),
			expectedBook: &Book{
				Id:     1,
				Title:  "book",
				Author: "author",
				Year:   2029,
				Size:   300,
				Rate:   4.33,
			},
		},
		{
			name:         "json пустой книги",
			bookBytes:    []byte("{}"),
			expectedBook: &Book{},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			book := &Book{}
			err := book.UnmarshalJSON(testCase.bookBytes)

			require.NoError(t, err)
			require.Equal(t, testCase.expectedBook, book)
		})
	}
}

func TestMarshalBooksJSON(t *testing.T) {
	testCases := []struct {
		name         string
		books        []*Book
		expectedJSON string
	}{
		{
			name:         "пустой слайс",
			books:        make([]*Book, 0),
			expectedJSON: "[]",
		},
		{
			name: "слайс с одним элементом",
			books: []*Book{
				{
					Id:     1,
					Title:  "book",
					Author: "author",
					Year:   2029,
					Size:   300,
					Rate:   4.33,
				},
			},
			expectedJSON: `[{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33}]`,
		},
		{
			name: "слайс с несколькими элементами",
			books: []*Book{
				{
					Id:     1,
					Title:  "book",
					Author: "author",
					Year:   2029,
					Size:   300,
					Rate:   4.33,
				}, {
					Id:     1,
					Title:  "book",
					Author: "author",
					Year:   2029,
					Size:   300,
					Rate:   4.33,
				},
			},
			expectedJSON: `[{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33},{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33}]`, //nolint:lll
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bytes, err := MarshalBooksJSON(testCase.books)

			require.NoError(t, err)
			require.Equal(t, testCase.expectedJSON, string(bytes))
		})
	}
}

func TestUnmarshalBooksJSON(t *testing.T) {
	testCases := []struct {
		name              string
		bytes             []byte
		expectedBookSlice []*Book
	}{
		{
			name:              "пустой слайс",
			bytes:             []byte("[]"),
			expectedBookSlice: make([]*Book, 0),
		},
		{
			name:  "слайс с одним элементом",
			bytes: []byte(`[{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33}]`),
			expectedBookSlice: []*Book{
				{
					Id:     1,
					Title:  "book",
					Author: "author",
					Year:   2029,
					Size:   300,
					Rate:   4.33,
				},
			},
		},
		{
			name: "слайс с несколькими элементами",
			bytes: []byte(`[
								{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33},
								{"id":1,"title":"book","author":"author","year":2029,"size":300,"rate":4.33}
							]`),
			expectedBookSlice: []*Book{
				{
					Id:     1,
					Title:  "book",
					Author: "author",
					Year:   2029,
					Size:   300,
					Rate:   4.33,
				}, {
					Id:     1,
					Title:  "book",
					Author: "author",
					Year:   2029,
					Size:   300,
					Rate:   4.33,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			books, err := UnmarshalBooksJSON(testCase.bytes)

			require.NoError(t, err)
			require.Equal(t, testCase.expectedBookSlice, books)
		})
	}
}
