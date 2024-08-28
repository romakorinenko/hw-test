package fixapp

import (
	"errors"
	"testing"

	"github.com/fixme_my_friend/hw06_testing/fixapp/model"
	"github.com/stretchr/testify/require"
)

func TestPrintStaffFromFile(t *testing.T) {
	testCases := []struct {
		name          string
		fileName      string
		expectedStaff []model.Employee
		expectedError error
	}{
		{
			name:     "передаем data.json",
			fileName: "data.json",
			expectedStaff: []model.Employee{
				{UserID: 10, Age: 25, Name: "Rob", DepartmentID: 3},
				{UserID: 11, Age: 30, Name: "George", DepartmentID: 2},
			},
			expectedError: nil,
		},
		{
			name:     "передаем пустую строку как fileName",
			fileName: "",
			expectedStaff: []model.Employee{
				{UserID: 10, Age: 25, Name: "Rob", DepartmentID: 3},
				{UserID: 11, Age: 30, Name: "George", DepartmentID: 2},
			},
			expectedError: nil,
		},
		{
			name:          "передаем несуществующий файл",
			fileName:      "test.json",
			expectedStaff: nil,
			expectedError: errors.New("cannot read staff from test.json"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			staff, err := readAndPrint(testCase.fileName)
			require.Equal(t, testCase.expectedStaff, staff)
			require.Equal(t, testCase.expectedError, err)
		})
	}
}
