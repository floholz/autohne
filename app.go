package main

import (
	. "autohne/src"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
)

var twitch = MakeTwitchApi()
var videoUtils = NewVideoUtils(true)
var youtube = YoutubeApi{}

func main() {
	var youtubeEnabled bool
	pflag.BoolVarP(&youtubeEnabled, "youtube", "y", false, "Use YouTube")
	var tiktokEnabled bool
	pflag.BoolVarP(&tiktokEnabled, "tiktok", "t", false, "Use TikTok")
	var instagramEnabled bool
	pflag.BoolVarP(&instagramEnabled, "instagram", "i", false, "Use Instagram")
	pflag.Parse()

	command := pflag.Arg(0)

	switch command {
	case "download":
		downloadClips()
	case "create":
		createShort()
	case "upload":
		uploadShort(youtubeEnabled, tiktokEnabled, instagramEnabled)
	}
}

func uploadShort(yt bool, tt bool, ig bool) {
	fmt.Printf("Upload to:\n\tYouTube: %t\n\tTikTok: %t\n\tInstagram: %t\n", yt, tt, ig)

	file, err := os.ReadFile("./assets/.ignore/short.mp4")
	if err != nil {
		log.Fatal(err)
	}

	if yt {
		uploadToYouTube(file)
	}
	if tt {
		// do TikTok stuff
	}
	if ig {
		// do Instagram stuff
	}
}

func uploadToYouTube(file []byte) {

	videoData := NewYoutubeVideData(
		"Souvenir Dragon Lore owner btw #ohnepixel",
		"Clip from ohnepixel stream",
		"20",
		"ohnepixel, clip, twitch, stream, highlight",
		"public",
	)
	youtube.UploadVideo(file, videoData)
}

func createShort() {
	file, err := os.ReadFile("./assets/.ignore/out/clip_souvenir_dragonlore_owner_btw_.mp4")
	if err != nil {
		log.Fatal(err)
	}

	// short := videoUtils.CreateShort(file)
	short := videoUtils.CreateShortFromFullVid(file)

	err = os.WriteFile("assets/.ignore/short.mp4", short, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadClips() {
	clips := twitch.GetClips()

	jayson, _ := json.MarshalIndent(clips, "", "  ")
	fmt.Println(string(jayson))

	for _, clip := range clips {
		clip.DownloadClip()
	}
}
