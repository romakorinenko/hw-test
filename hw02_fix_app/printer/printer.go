package printer

import (
	"fmt"

	"github.com/hw-test/hw02_fix_app/model"
)

func PrintStaff(staff []model.Employee) {
	for _, employee := range staff {
		fmt.Println(employee.String())
	}
}
