package fileactions

import (
	"testing"
	"time"

	"github.com/davidparks11/file-renamer/pkg/config"
	"github.com/davidparks11/file-renamer/pkg/fileretriever"
	"github.com/davidparks11/file-renamer/pkg/fileretriever/fileretrieveriface"
	"github.com/davidparks11/file-renamer/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFileActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FileActions Suite")

	var _ = Describe("FileActions", func() {
		mockRetriever := &fileretriever.MockFileRetriever{}
		mockToday, _ := time.Parse(time.RFC3339, "2020-08-24T19:33:44.561Z")
		fileAction := &Renamer{
			logger:        &logger.MockLogger{},
			fileRetriever: mockRetriever,
			name:          "testfileaction",
			config: &config.Config{
				PersistentWords: []string{"foo", "bar"},
				NameDelimiter:   "_",
			},
			nameFlushDate:  mockToday.Add(-72 * time.Hour),
			processedFiles: map[string]bool{"2010_0111_0.mov": true},
		}
		Describe("generateNewName()", func() {
			It("Should return a generated name with an incremented duplicate number", func() {
				actual, err := fileAction.generateNewName("foogetsdeletedbar.mov", "2020-10-19T04:49:06.334Z")
				expected := "foo_bar_2020_1019_1.mov"
				Expect(err).To(BeNil())
				Expect(fileAction.processedFiles["foo_bar_2020_1019_0.mov"]).To(Equal(true))
				Expect(actual).To(Equal(expected))
			})
		})

		Describe("parseRFC3339()", func() {
			It("should return nothing when given an invalid timestamp", func() {
				actual, err := fileAction.parseRFC3339("2020-03-24 04:45:20")
				expected := ""
				Expect(actual).To(Equal(expected))
				Expect(err).ToNot(BeNil())
			})
			It("should return the format YYYY_MMDD when given a valid timestamp", func() {
				actual, err := fileAction.parseRFC3339("2020-08-31T19:33:44.561Z")
				expected := "2020_0831"
				Expect(actual).To(Equal(expected))
				Expect(err).To(BeNil())
			})
		})

		Describe("Run()", func() {
			mockFileInfo := []*fileretrieveriface.RenameInfo{
				{
					ID:          "11111111",
					Name:        "foofile1.mov",
					CreatedDate: "2020-08-31T19:33:44.561Z",
				},
				{
					ID:          "22222222",
					Name:        "foofile1.mov",
					CreatedDate: "2020-08-31T19:33:44.561Z",
				},
				{
					ID:          "33333333",
					Name:        "foobarfile2.mov",
					CreatedDate: "2020-08-31T17:33:44.561Z",
				},
				{
					ID:          "44444444",
					Name:        "file3.mov",
					CreatedDate: "2020-08-28T19:33:44.561Z",
				},
			}
			mockRetriever.On("GetFileInfo").Return(mockFileInfo, nil)

			processedFiles := map[string]bool{"2020_0828_0.mov": true}

			now = func() time.Time { return mockToday }

			mockRetriever.On("GetProcessedFiles", mockToday.Add(time.Hour * 48).Format(time.RFC3339)).Return(processedFiles)

			updatedFiles := []*fileretrieveriface.RenameInfo {
				{ID:"11111111", Name:"foo_2020_0831_0.mov", CreatedDate:"2020-08-31T19:33:44.561Z"},
				{ID:"22222222", Name:"foo_2020_0831_1.mov", CreatedDate:"2020-08-31T19:33:44.561Z"},
				{ID:"33333333", Name:"foo_bar_2020_0831_0.mov", CreatedDate:"2020-08-31T17:33:44.561Z"},
				{ID:"44444444", Name:"2020_0828_1.mov", CreatedDate:"2020-08-28T19:33:44.561Z"},				
			}

			mockRetriever.On("UpdateFile", updatedFiles[0]).Return(nil)
			mockRetriever.On("UpdateFile", updatedFiles[1]).Return(nil)
			mockRetriever.On("UpdateFile", updatedFiles[2]).Return(nil)
			mockRetriever.On("UpdateFile", updatedFiles[3]).Return(nil)

			fileAction.Run()

			mockRetriever.AssertExpectations(t)

			//Check that processedFiles were updated
			Expect(fileAction.processedFiles).NotTo(Equal("2010_0111_0.mov"))
		})
	})

}
