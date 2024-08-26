package fixapp

import (
	"github.com/fixme_my_friend/hw06_testing/fixapp/model"
	"github.com/fixme_my_friend/hw06_testing/fixapp/reader"
)

func GetStaffFromFile(filePath string) ([]model.Employee, error) {
	if len(filePath) == 0 {
		filePath = "data.json"
	}

	staff, err := reader.ReadJSON(filePath)
	if err != nil {
		return nil, err
	}

	return staff, nil
}
