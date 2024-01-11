package src

import (
	"log"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type VideoUtils struct {
}

func (vu VideoUtils) AddWatermark() {
	overlay := ffmpeg.Input("./assets/logo_wm.png").Filter("scale", ffmpeg.Args{"-1:48"})
	err := ffmpeg.Filter(
		[]*ffmpeg.Stream{
			ffmpeg.Input("./assets/.private/clip.mp4"),
			overlay,
		}, "overlay", ffmpeg.Args{"(abs(main_h/2)-overlay_w-10):(main_h-overlay_h-10)"}).
		Filter("crop", ffmpeg.Args{"abs(ih/2):ih:0:0"}).
		Output("./assets/.private/watermark.mp4",
			ffmpeg.KwArgs{"map": "0:a"},
			ffmpeg.KwArgs{"t": 7}).
		OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		log.Fatal(err)
	}
}

func (vu VideoUtils) OpenCV() {

}
