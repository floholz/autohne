package main

import (
	. "autohne/src"
	"encoding/json"
	"fmt"
	"time"
)

var twitch = MakeTwitchApi()
var videoUtils = VideoUtils{}

func main() {
	// run()
	// tick()
	videoUtils.AddWatermark()
}

func tick() {
	clips := twitch.GetClips()

	jayson, _ := json.MarshalIndent(clips, "", "  ")
	fmt.Println(string(jayson))

	for _, clip := range clips {
		clip.DownloadClip()
	}
}

func finally() {
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
