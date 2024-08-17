package main

import (
	"fmt"
	"github.com/fixme_my_friend/hw02_fix_app/printer"
	"github.com/fixme_my_friend/hw02_fix_app/reader"
	"log/slog"
)

func main() {
	var path string
	fmt.Printf("Enter data file path: ")
	if n, err := fmt.Scanln(&path); err != nil || n != 1 {
		slog.Error("Cannot scan file name from command line", slog.Any("error", err))
	}

	if len(path) == 0 {
		path = "data.json"
	}

	staff, err := reader.ReadJSON(path)
	if err != nil {
		slog.Error("Cannot read from file", slog.Any("error", err))
	}

	printer.PrintStaff(staff)
}
