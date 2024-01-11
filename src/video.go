package src

import (
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoUtils struct {
}

func (vu VideoUtils) AddWatermark() {
	// show watermark with size 64:-1 in the top left corner after seconds 1
	overlay := ffmpeg.Input("./assets/rocket_1f680.png").Filter("scale", ffmpeg.Args{"128:-1"})
	err := ffmpeg.Filter(
		[]*ffmpeg.Stream{
			ffmpeg.Input("./assets/clip.mp4"),
			overlay,
		}, "overlay", ffmpeg.Args{"10:10"}).
		Output("./assets/watermark.mp4", ffmpeg.KwArgs{"map": "0:a"}, ffmpeg.KwArgs{"t": 5}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Fatal(err)
	}
}
