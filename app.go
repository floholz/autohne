package main

import (
	. "autohne/src"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var twitch = MakeTwitchApi()
var videoUtils = NewVideoUtils(true)
var youtube = YoutubeApi{}

func main() {
	//downloadClips()
	//createShort()
	//uploadShort()
}

func uploadShort() {
	file, err := os.ReadFile("./assets/.ignore/short.mp4")
	if err != nil {
		log.Fatal(err)
	}

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

func run() {
	done := make(chan bool)
	ticker := time.NewTicker(time.Second * 5)
	defer finally()

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Done !!")
				ticker.Stop()
				return
			case <-ticker.C:
				fmt.Printf("%s | still alive!\n", time.Now().Format(time.RFC3339))
				tick()
			}
		}
	}()

	tick()
	// run for 15 seconds
	time.Sleep(15 * time.Second)
	done <- true
}
func tick() {

}
func finally() {
}
