package fileretriever

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/davidparks11/file-renamer/pkg/config"
	"github.com/davidparks11/file-renamer/pkg/fileretriever/fileretrieveriface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	"google.golang.org/api/option"
)

var _ fileretrieveriface.FileRetriever = &FileRetriever{}

const (
	MAX_FILE_RESULTS = 10
	PROCESSED_PROP_FIELD = "processed"
)

//FileRetriever type that serves to get files and write changes to files
type FileRetriever struct {
	logger loggeriface.Service
	config *config.Config
	drive  *drive.Service
}

//GetFileInfo returns all files that match description from config
func (d *FileRetriever) GetFileInfo() (*[]os.FileInfo, error) {
	files, err := d.drive.Files.List().
	DriveId(d.config.ParentDirID).
	Q("mimeType='application/vnd.google-apps.video'" +
		"and " + PROCESSED_PROP_FIELD + " = false" +
		"and trashed = false").
	MaxResults(MAX_FILE_RESULTS).Do() 
	if err != nil {
		d.logger.Error("Failed to get files: " + err.Error())
	}
	for _, v := range files.Items {
		d.logger.Info(v.TeamDriveId)
	}

	return nil, nil
}

func (d *FileRetriever) UpdateFiles() error {
	return nil
}

//All Code below was edited but used from the Google drive quickstart quide for golang at https://developers.google.com/drive/api/v3/quickstart/go

//NewFileRetriever serves a file retriever
func NewFileRetriever(logger loggeriface.Service, config *config.Config) fileretrieveriface.FileRetriever {
	fileRetriever := &FileRetriever{
		logger: logger,
		config: config,
	}

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
func (d *FileRetriever) getClient(oauthConfig *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	token, err := d.tokenFromFile()
	if err != nil {
		token = d.getTokenFromWeb(oauthConfig)
		d.saveToken(token)
	}
	return oauthConfig.Client(context.Background(), token)
}

// Request a token from the web, then returns the retrieved token.
func (d *FileRetriever) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	d.logger.Info("Go to the following link in your browser then type the "+
		"authorization code: " + authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		d.logger.Fatal("Unable to read authorization code: " + err.Error())
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		d.logger.Fatal("Unable to retrieve token from web: " + err.Error())
	}
	return tok
}

// Retrieves a token from a local file.
func (d *FileRetriever) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(d.config.TokenPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func (d *FileRetriever) saveToken(token *oauth2.Token) {
	d.logger.Info("Saving credential file to: " + d.config.TokenPath)
	f, err := os.OpenFile(d.config.TokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		d.logger.Fatal("Unable to cache oauth token:: " + err.Error())
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
