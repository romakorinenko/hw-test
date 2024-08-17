package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

import "github.com/fixme_my_friend/hw02_fix_app/model"

func ReadJSON(filePath string) ([]model.Employee, error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return nil, err
	}

	var data []model.Employee
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
