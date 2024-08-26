package fixapp

import (
	"testing"

	"github.com/fixme_my_friend/hw06_testing/fixapp/model"
	"github.com/stretchr/testify/require"
)

func TestGetStaffFromFile(t *testing.T) {
	testCases := []struct {
		name           string
		fileName       string
		expectedResult []model.Employee
		containsError  bool
	}{
		{
			name:     "Должен вернуть двух сотрудников без ошибок",
			fileName: "data.json",
			expectedResult: []model.Employee{
				{UserID: 10, Age: 25, Name: "Rob", DepartmentID: 3},
				{UserID: 11, Age: 30, Name: "George", DepartmentID: 2},
			},
			containsError: false,
		},
		{
			name:           "Должен вернуть ошибку",
			fileName:       "wrong.json",
			expectedResult: nil,
			containsError:  true,
		},
		{
			name:     "Передаем строку с пустым именем файла",
			fileName: "",
			expectedResult: []model.Employee{
				{UserID: 10, Age: 25, Name: "Rob", DepartmentID: 3},
				{UserID: 11, Age: 30, Name: "George", DepartmentID: 2},
			},
			containsError: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			staff, err := GetStaffFromFile(testCase.fileName)
			require.Equal(t, testCase.expectedResult, staff)
			require.Equal(t, testCase.containsError, err != nil)
		})
	}
}
