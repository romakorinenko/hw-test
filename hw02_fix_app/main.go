package main

import (
	"fmt"
	"github.com/fixme_my_friend/hw02_fix_app/printer"
	"github.com/fixme_my_friend/hw02_fix_app/reader"
	"log/slog"
)

func main() {
	var path string

	fmt.Print("Enter data file path: ")
	_, _ = fmt.Scanln(&path)

	if len(path) == 0 {
		path = "data.json"
	}

	staff, err := reader.ReadJSON(path)
	if err != nil {
		slog.Error("Cannot read from file", slog.Any("error", err))
	}

	printer.PrintStaff(staff)
}
