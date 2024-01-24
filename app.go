package main

import (
	. "autohne/src"
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
)

var appConfig = NewAppConfig()

var twitch = NewTwitchApi(appConfig.Download.TwitchConfig)
var videoUtils = NewVideoUtils(true)
var youtube = NewYoutubeApi(appConfig.Upload.YoutubeConfig)

func main() {
	var cmdDownload bool
	pflag.BoolVarP(&cmdDownload, "download", "d", false, "Download newest clips")
	var cmdCreate bool
	pflag.BoolVarP(&cmdCreate, "create", "c", false, "Create short format videos from clips")
	var cmdUpload bool
	pflag.BoolVarP(&cmdUpload, "upload", "u", false, "Upload video to the specified platforms")

	var youtubeEnabled bool
	pflag.BoolVarP(&youtubeEnabled, "Youtube", "Y", false, "Use YouTube")
	var tiktokEnabled bool
	pflag.BoolVarP(&tiktokEnabled, "Tiktok", "T", false, "Use TikTok")
	var instagramEnabled bool
	pflag.BoolVarP(&instagramEnabled, "Instagram", "I", false, "Use Instagram")

	var debug bool
	pflag.BoolVar(&debug, "debug", false, "Debug mode")
	var help bool
	pflag.BoolVarP(&help, "help", "h", false, "Display command options and flags")

	pflag.Parse()

	if help {
		pflag.CommandLine.SortFlags = false
		pflag.PrintDefaults()
		os.Exit(0)
	}

	if cmdDownload {
		downloadClips()
	}
	if cmdCreate {
		createShort()
	}
	if cmdUpload {
		uploadShort(youtubeEnabled, tiktokEnabled, instagramEnabled)
	} else {
		if youtubeEnabled || tiktokEnabled || instagramEnabled {
			err := fmt.Errorf("'-Y', '-T' and '-I' do only work in combination with '-u'")
			fmt.Println(err.Error())
		}
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
