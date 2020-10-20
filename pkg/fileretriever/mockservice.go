package fileretriever

import (
	"fmt"

	"github.com/davidparks11/file-renamer/pkg/fileretriever/fileretrieveriface"
	"github.com/stretchr/testify/mock"
)


var _ fileretrieveriface.FileRetriever = &MockFileRetriever{}

type MockFileRetriever struct {
	mock.Mock
}

//GetFileInfo mocks a call to gdrive
func (m *MockFileRetriever)GetFileInfo() ([]*fileretrieveriface.RenameInfo, error) {
	// mockFileInfo := []*fileretrieveriface.RenameInfo{
	// 	{
	// 		ID:          "11111111",
	// 		Name:        "foofile1.mov",
	// 		CreatedDate: "2020-08-31T19:33:44.561Z",
	// 	},
	// 	{
	// 		ID:          "22222222",
	// 		Name:        "foofile1.mov",
	// 		CreatedDate: "2020-08-31T19:33:44.561Z",
	// 	},
	// 	{
	// 		ID:          "33333333",
	// 		Name:        "foobarfile2.mov",
	// 		CreatedDate: "2020-08-31T17:33:44.561Z",
	// 	},
	// 	{
	// 		ID:          "44444444",
	// 		Name:        "file3.mov",
	// 		CreatedDate: "2020-08-28T19:33:44.561Z",
	// 	},
	// }
	args := m.Called()
	return args.Get(0).([]*fileretrieveriface.RenameInfo), nil
}

//CheckForDuplicate mocks a call to gdrive to check if a file exists
func (m *MockFileRetriever) IsUniqueName(name string) bool {
args := m.Called(name)
return args.Bool(0)
}

//UpdateFile it just returns nil
func (m *MockFileRetriever) UpdateFile(info *fileretrieveriface.RenameInfo) error {
	fmt.Println("\n\n\n\n\n~~~~~~~~~~~~~~~"+info.Name+"~~~~~~~~~~~~~\n\n\n\n\n")
	args := m.Called(info)
	return args.Error(0)
}

func (m *MockFileRetriever) GetProcessedFiles(date string) map[string]bool {
	//files := map[string]bool {"foo_2020_0901.mov": true}
	args := m.Called(date)
	return args.Get(0).(map[string]bool)
}
