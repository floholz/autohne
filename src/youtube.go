package src

import (
	"bytes"
	"fmt"
	"google.golang.org/api/youtube/v3"
	"log"
	"strings"
)

type YoutubeApi struct {
}

type YoutubeVideoData struct {
	title       string
	description string
	category    string
	keywords    string
	privacy     string
}

func NewYoutubeVideData(title string, description string, category string, keywords string, privacy string) YoutubeVideoData {
	vd := YoutubeVideoData{}
	vd.title = title
	vd.description = description
	vd.category = category
	vd.keywords = keywords
	vd.privacy = privacy
	return vd
}

func (yt *YoutubeApi) UploadVideo(video []byte, videoData YoutubeVideoData) {

	fileBuf := bytes.NewBuffer(video)

	client := getClient(youtube.YoutubeUploadScope)

	service, err := youtube.New(client)
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	upload := &youtube.Video{
		Snippet: &youtube.VideoSnippet{
			Title:       videoData.title,
			Description: videoData.description,
			CategoryId:  videoData.category,
		},
		Status: &youtube.VideoStatus{PrivacyStatus: videoData.privacy},
	}

	// The API returns a 400 Bad Request response if tags is an empty string.
	if strings.Trim(videoData.keywords, "") != "" {
		upload.Snippet.Tags = strings.Split(videoData.keywords, ",")
	}

	part := []string{"snippet", "status"}
	call := service.Videos.Insert(part, upload)

	response, err := call.Media(fileBuf).Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Upload successful! Video ID: %v\n", response.Id)
}
