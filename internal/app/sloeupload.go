package app

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sloe/gosloe/internal/config"
	"github.com/sloe/gosloe/internal/domain"
)

var translateEvent = map[string]string{
	"townbumps2019": "Cambridge Town Bumps 2019",
	"townbumps2021": "Cambridge Town Bumps 2021",
}

var translateDay = map[string]string{
	"sun":   "Sunday",
	"mon":   "Monday",
	"tues":  "Tuesday",
	"wed":   "Wednesday",
	"thurs": "Thursday",
	"fri":   "Friday",
	"sat":   "Saturday",
}

func doMuseUpload(filePath, title string, config config.Config) error {
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	titleSubbed := strings.Replace(title, " ", "_", -1)
	part, _ := writer.CreateFormFile("file", titleSubbed)
	io.Copy(part, file)
	writer.Close()

	r, _ := http.NewRequest("POST", config.MuseUploadUrl, body)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	r.Header.Add("Key", config.MuseApiKey)
	client := &http.Client{}
	uploadResponse, err := client.Do(r)
	log.WithError(err).WithFields(log.Fields{"uploadResponse": uploadResponse}).Info("Uploaded file")
	return err
}

func SloeUpload(tree domain.SloeTree, config config.Config) error {

	// Hardcoded for now

	for _, item := range tree.Items {
		if strings.HasPrefix(item.SubTree, "final/derived/townbumps") {
			tagFilePath := item.Location + ".done"
			if _, err := os.Stat(tagFilePath); os.IsNotExist(err) {
				videoFilePath := filepath.Join(filepath.Dir(item.Location), item.Leafname)

				titleRegexp := regexp.MustCompile(`/([^/]+)/([^/]+)/([^/]+)/([^/]+)`)
				match := titleRegexp.FindStringSubmatch(item.SubTree)
				if len(match) == 0 {
					log.WithFields(log.Fields{"SubTree": item.SubTree}).Error("Could not decode subtree")
				} else {
					event, ok := translateEvent[match[2]]
					if ok {

						day, ok := translateDay[match[3]]
						if !ok {
							day = ""
						}
						title := strings.Replace(item.Leafname, " ytf", " normal speed", 1)
						title = strings.Replace(title, " yt8", " slow motion", 1)
						title = strings.Replace(title, " yt4", " slow motion", 1)
						title = event + " " + day + " " + title

						log.WithFields(log.Fields{"Location": item.Location, "Title": title}).Info("Uploading file")
						err = doMuseUpload(videoFilePath, title, config)
						if err != nil {
							log.WithError(err).WithFields(log.Fields{"Location": item.Location, "Title": title}).Error("Failed to upload file")
						} else {
							file, err := os.Create(tagFilePath)
							if err != nil {
								log.WithError(err).WithFields(log.Fields{"TagFilePath": tagFilePath}).Error("Failed to create tag file")
							} else {
								_ = file.Close()
							}
						}
					}
				}
			}

		}
	}

	return nil
}
