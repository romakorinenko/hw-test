package reader

import (
	"encoding/json"
	"io"
	"os"

	"github.com/fixme_my_friend/hw06_testing/fixapp/model"
)

func ReadJSON(filePath string) ([]model.Employee, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var data []model.Employee
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
