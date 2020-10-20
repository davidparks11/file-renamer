package fileretriever

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/davidparks11/file-renamer/pkg/config"
	"github.com/davidparks11/file-renamer/pkg/fileretriever/fileretrieveriface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/time/rate"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

const (
	fileProcessedFlag = "file-renamer-processed"
	readLimitPerMinute = 300
	writeLimitPerMinute = 60
)

var _ fileretrieveriface.FileRetriever = &FileRetriever{}

//FileRetriever type that serves to get files and write changes to files
type FileRetriever struct {
	logger loggeriface.Service
	config *config.Config
	drive  *drive.Service
	queryableFolders []string
	readLimiter *rate.Limiter
	writeLimiter *rate.Limiter
}

func (f *FileRetriever) getSubFolders(parentFolder string) ([]string, error) {
	//slice to hold parent dir and all children dirs under it
	folderIds := []string{parentFolder}
	requestCount := 0
	folderIndex := 0
	var query string
	for len(folderIds) != folderIndex {
		//Build query for files whose parents are in folderIds
		query = f.buildChildQuery(folderIds[folderIndex:])
		folderIndex = len(folderIds)

		//Set folder index to address the first of the next folder ids
		err := f.drive.Children.List(f.config.ParentDirID).
		MaxResults(1000).
		Q(query).
		Pages(
			context.TODO(),
			func(folders *drive.ChildList) error {
				requestCount++
				for _, v := range folders.Items {
					folderIds = append(folderIds, v.Id)
				}
				return nil
			},
		)
		
		if err != nil {
			return nil, err
		}
	}

	return folderIds, nil
}

func (f *FileRetriever) getFilesFromFolders(folderIds []string) ([]*fileretrieveriface.RenameInfo, error) {
	//After all child folder of the config-parent dir have been found
	//query for any files to rename 
	query := f.buildFileQuery(folderIds)
	//only get files that have not been processed
	query += "and not (" + processedQuery + ")"
	var files []*fileretrieveriface.RenameInfo
	err := f.drive.Files.List().
		MaxResults(1000).
		Q(query).
		Pages(
			context.TODO(),
			func(fileList *drive.FileList) error {
				for _, v := range fileList.Items {
					
					files = append(files, &fileretrieveriface.RenameInfo{
						ID: v.Id,
						Name: v.Title,
						CreatedDate: v.CreatedDate,
					})
					
				}
				return nil
			},
		)

	if err != nil {
		f.logger.Error(err.Error())
	}
	if len(files) == 0 {
		f.logger.Info("Couldn't find any files")
	} else {
		f.logger.Info(fmt.Sprintf("Found %d files", len(files)))
	}	

	return files, nil
}

//GetFileInfo returns all files that match description from config
func (f *FileRetriever) GetFileInfo() ([]*fileretrieveriface.RenameInfo, error) {
	folderIds, err := f.getSubFolders(f.config.ParentDirID)
	if err != nil {
		return nil, err
	}
	//retain queryable folders for duplicate search
	f.queryableFolders = folderIds
	return f.getFilesFromFolders(folderIds)
}

func (f *FileRetriever) buildChildQuery(folderIds []string) (query string) {
	if len(folderIds) == 0 {
		return ""
	}
	query = fmt.Sprintf("'%s' in parents ", folderIds[0]) 
	for i := 1; i < len(folderIds); i++ {
		query += fmt.Sprintf("or '%s' in parents ", folderIds[i])
	}
	query += "and mimeType = 'application/vnd.google-apps.folder'"
	return query
}

func (f *FileRetriever) buildFileQuery(folderIds []string) (query string) {
	for i := 0; i < len(folderIds); i++ {
		if i == 0 {
			query += "("
		} else {
			query += "or "
		}
		query += fmt.Sprintf("'%s' in parents ", folderIds[i])
		if i == len(folderIds) - 1 {
			query += ") "
		}
	}
	
	for i, v := range f.config.FileExtensions {
		if i == 0 {
			query += "and ("
		} else {
			query += "or "
		}
		query += fmt.Sprintf(" title contains '.%s' ", v)
		if i == len(f.config.FileExtensions) - 1 {
			query += ") "
		}
	}
	return query
}

const processedQuery = "properties has {key='" + fileProcessedFlag + "' and value='true' and visibility='PUBLIC'}"

func (f *FileRetriever) GetProcessedFiles(date string) map[string]bool {
	var processedFiles map[string]bool
	query := f.buildFileQuery(f.queryableFolders)
	query += "and " + processedQuery
	err := f.drive.Files.List().
		MaxResults(1000).
		Q(query).
		Pages(
			context.TODO(),
			func(fileList *drive.FileList) error {
				for _, v := range fileList.Items {
					processedFiles[v.Title] = true
				}
				return nil
			},
		)

	if err != nil {
		f.logger.Error(err.Error())
	}
	return processedFiles
}

//IsUniqueName returns nil if no file is found with the same name, otherwise, returns false
func (f *FileRetriever) IsUniqueName(name string) bool {
	response, err := f.drive.Files.List().
		Q("title = '" + name + "'").Do()

	if err != nil {
		f.logger.Error(err.Error())
	}
	if len(response.Items) > 0 {
		return false
	}
	return true
}

//UpdateFile gives the file a new name and sets a custom property to true on the file
func (f *FileRetriever) UpdateFile(info *fileretrieveriface.RenameInfo) error {
	processedProp := &drive.Property{
		Key: fileProcessedFlag,
		Value: "true",
		Visibility: "PUBLIC",
	}

	file := &drive.File {
		Title: info.Name,
		Properties: []*drive.Property{processedProp},
	}
	_, err := f.drive.Files.Update(info.ID, file).Do()
	return err
}

//All Code below was edited but used from the Google drive quickstart quide for golang at https://developers.google.com/drive/api/v3/quickstart/go

//NewFileRetriever serves a file retriever
func NewFileRetriever(logger loggeriface.Service, config *config.Config) fileretrieveriface.FileRetriever {
	fileRetriever := &FileRetriever{
		logger: logger,
		config: config,
	}

	//establish rate limiter for reads/writes
	wl := rate.NewLimiter(rate.Every(time.Minute), writeLimitPerMinute)
	rl := rate.NewLimiter(rate.Every(time.Minute), readLimitPerMinute)
	fileRetriever.readLimiter = rl
	fileRetriever.writeLimiter = wl
	
	b, err := ioutil.ReadFile(config.CredentialsPath)
	if err != nil {
		logger.Fatal("Unable to read client secret file: " + err.Error())
	}

	// If modifying these scopes, delete your previously saved token.json.
	oauthConfig, err := google.ConfigFromJSON(b, drive.DriveScope)
	if err != nil {
		logger.Fatal("Unable to parse client secret file to config: " + err.Error())
	}
	client := fileRetriever.getClient(oauthConfig)
	
	drive, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logger.Fatal("Unable to retrieve Drive client: " + err.Error())
	}

	fileRetriever.drive = drive
	return fileRetriever
}

// Retrieve a token, saves the token, then returns the generated client.
func (f *FileRetriever) getClient(oauthConfig *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	token, err := f.tokenFromFile()
	if err != nil {
		token = f.getTokenFromWeb(oauthConfig)
		f.saveToken(token)
	}
	return oauthConfig.Client(context.Background(), token)
}

// Request a token from the web, then returns the retrieved token.
func (f *FileRetriever) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	f.logger.Info("Go to the following link in your browser then type the "+
		"authorization code: " + authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		f.logger.Fatal("Unable to read authorization code: " + err.Error())
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		f.logger.Fatal("Unable to retrieve token from web: " + err.Error())
	}
	return tok
}

// Retrieves a token from a local file.
func (f *FileRetriever) tokenFromFile() (*oauth2.Token, error) {
	file, err := os.Open(f.config.TokenPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (f *FileRetriever) saveToken(token *oauth2.Token) {
	f.logger.Info("Saving credential file to: " + f.config.TokenPath)
	file, err := os.OpenFile(f.config.TokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		f.logger.Fatal("Unable to cache oauth token:: " + err.Error())
	}
	defer file.Close()
	json.NewEncoder(file).Encode(token)
}
