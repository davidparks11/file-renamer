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

type node struct {
	id string
	name string
	children []*node
}

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

func (d *FileRetriever) buildFileTree() {
	//get map of parent to file
	var files map[string]drive.File
	err := d.drive.Files.
		List().
		MaxResults(3000).
		Pages(
			context.TODO(),
			func(folders *drive.FileList) error {
				for _, v := range folders.Items {
					files[v.Parents[0].Id] = *v
				}
				return nil
			},
		)

	if err != nil {
		d.logger.Error(err.Error())
	}

}

//GetFileInfo returns all files that match description from config
func (d *FileRetriever) GetFileInfo() (*[]os.FileInfo, error) {
	// please, err := d.drive.Files.List().Do()
	// if err != nil {
	// 	return nil, err
	// }
	// if len(please.Items) == 0 {
	// 	d.logger.Info("No Drives available")
	// } else {
	// 	for _, v := range please.Items {
	// 		d.logger.Info("Name: " + v.Title + "\tID: " + v.Id)
	// 	}
	// } mimeType = 'application/vnd.google-apps.folder'

	
	var folderIds []string
	err := d.drive.Children.List(d.config.ParentDirID).
		MaxResults(200).
		Pages(
			context.TODO(),
			func(folders *drive.ChildList) error {
				for _, v := range folders.Items {
					folderIds = append(folderIds, v.Id)
				}
				return nil
			},
		)

	if err != nil {
		d.logger.Error(err.Error())
	}
	if len(folderIds) == 0 {
		d.logger.Info("Couldn't find any IDs")
	} else {
		d.logger.Info(fmt.Sprintf("Found %d ids", len(folderIds)))
	}

	query := d.buildFileQuery()
	d.logger.Info("query: " + query)

	var videos []*drive.File
	err = d.drive.Files.List().
		MaxResults(200).
		Q(query).
		Pages(
			context.TODO(),
			func(videoList *drive.FileList) error {
				for _, v := range videoList.Items {
					videos = append(videos, v)
				}
				return nil
			},
		)
	
	if err != nil {
		d.logger.Error(err.Error())
	}
	if len(videos) == 0 {
		d.logger.Info("Couldn't find any videos")
	} else {
		d.logger.Info(fmt.Sprintf("Found %d videos", len(videos)))
	}

	for i, v := range videos {
		d.logger.Info(fmt.Sprintf("Index: %d\t%s", i, v.Title))
	}



	return nil, nil
}

func (d *FileRetriever) buildFileQuery() string {
	query := fmt.Sprintf("'%s' in parents ", d.config.ParentDirID) 
	// for _, v := range d.config.FileExtensions {
	// 	query += fmt.Sprintf("or name contains '.%s' ", v)
	// }
	return query
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
