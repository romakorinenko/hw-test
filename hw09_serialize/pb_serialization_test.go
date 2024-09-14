package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestSerializeAndDeserializeBooksPB(t *testing.T) {
	testCases := []struct {
		name  string
		books *Books
	}{
		{
			name: "Books содержит пустой слайс",
			books: &Books{
				Books: []*Book{},
			},
		},
		{
			name: "Books содержит слайс с одним элементом",
			books: &Books{
				Books: []*Book{
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
		},
		{
			name: "Books содержит слайс с несколькими элементами",
			books: &Books{
				Books: []*Book{
					{
						Id:     1,
						Title:  "book",
						Author: "author",
						Year:   2029,
						Size:   300,
						Rate:   4.33,
					},
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
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			bytes, err := SerializeBooksPB(testCase.books)
			require.NoError(t, err)

			books := &Books{}
			err = DeserializeBooksPB(books, bytes)
			require.NoError(t, err)
			require.True(t, proto.Equal(testCase.books, books))
		})
	}
}
