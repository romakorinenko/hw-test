package printer

import (
	"fmt"

	"github.com/fixme_my_friend/hw02_fix_app/model"
)

func PrintStaff(staff []model.Employee) {
	for _, employee := range staff {
		fmt.Println(employee.String())
	}
}
