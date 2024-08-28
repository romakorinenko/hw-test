package printer

import (
	"fmt"

	"github.com/fixme_my_friend/hw06_testing/fixapp/model"
)

func PrintStaff(staff []model.Employee) {
	for _, employee := range staff {
		fmt.Println(employee.String())
	}
}
