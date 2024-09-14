package main

import "google.golang.org/protobuf/proto"

func SerializeBooksPB(books *Books) ([]byte, error) {
	return proto.Marshal(books)
}

func DeserializeBooksPB(books *Books, data []byte) error {
	return proto.Unmarshal(data, books)
}
