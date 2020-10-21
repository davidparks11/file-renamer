# File-Renamer
File-Renamer gathers files from google drive and renames using various configurations. One of the few non customizable portions of this name is the ending which takes the form "YYYYMMDD#.suffix", where the creation date populates for first portion, a duplicate number takes the place of '#'. a delimiter can be specified to separate various parts of the newly generated names. 

## Setup
For this project to run, the you must enable the google drive api and get OAuth credentials ([I have no idea what you're talking about](https://developers.google.com/drive/api/v3/quickstart/js)).  

## Configuration
There are a number of configuration options to aid in renaming files. It's recommended that one fills out the following list of configuration options.
- **cronSchedules**: Array of string representing schedules that the renamer will run on. For help on creating these schedules, visit [this help crontab website](https://crontab.guru/)
- **persistentWords**: Array of strings that will persist in titles in array index order when matched (case insensitive)  
- **parentDirID**: ID of the folder containing files that you want to rename. By traveling to the folder in google drive, you can find this ID in the in URL
- **nameDelimiter**: Single character to join all the file portions together (including persistent words). These would most commonly be "_" or "-"
- **fileExtensions**: Array of strings representing file extensions. Any file that has an extension in **fileExtensions** will be renamed
- **logLevel**: The granularity of which logs are recorded. **GIVEN AS A STRING**
    - 1: Errors only
    - 2: Warnings and errors
    - 3: information, warnings, and errors
- **logLocation**: Path for logs files to be written
- **credentialsPath**: Path to google drive credentials
- **tokenPath**: Path to google drive access token
- **runAtLaunch**: Runs the program on startup in addition to waiting for schedule. Mostly useful for testing.
- **logToConsole**: Logs to console in addition to log files
### EXAMPLE JSON 

```json
{
    "cronSchedules": ["30 12 * * 5"],
    "persistentWords": ["work", "portaits", "landscape"],
    "parentDirID": "folderIdHere",
    "nameDelimiter" : "_",
    "fileExtensions": ["pgp", "jpeg", "png"],
    "logLevel": "2",
    "logLocation": "C:\\Users\\youOrSomething\\logs",
    "credentialsPath": "resources\\credentials.json",
    "tokenPath": "resources\\token.json",
    "runAtLaunch": false,
    "logToConsole": false
}
```
## **NOTICE**
Upon running this program for the first time, if no token.json is found, then you will be prompt to visit a google site to proceed with the generation of your token. **logToConsole** must be set to **true** if you want your token.
## Usage
Assuming you have enabled the api, downloaded your credentials, and set your credentials path, you can run the following. 
```bash
go run cmd/main.go
```
To end the program, press ctrl+c.
## License
[MIT](https://choosealicense.com/licenses/mit/)
