package fixapp

import (
	"fmt"

	"github.com/fixme_my_friend/hw06_testing/fixapp/model"
	"github.com/fixme_my_friend/hw06_testing/fixapp/printer"
	"github.com/fixme_my_friend/hw06_testing/fixapp/reader"
)

func readAndPrint(path string) ([]model.Employee, error) {
	if len(path) == 0 {
		path = "data.json"
	}

	staff, err := reader.ReadJSON(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read staff from %v", path)
	}

	printer.PrintStaff(staff)

	return staff, nil
}
