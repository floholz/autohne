package src

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
)

type VideoUtils struct {
}

type FFmpegStream struct {
	*ffmpeg.Stream
}

func (vu VideoUtils) CreateShort(file []byte) []byte {
	inBuf := bytes.NewBuffer(file)
	outBuf := bytes.NewBuffer(nil)
	out := ffmpeg.Input("pipe:").Filter("scale", ffmpeg.Args{"-1:1080"})

	out = cropForShort(out)
	out = addOverlay(out, "assets/ohnepixel-clipper-mark.png", "500:-1", "25:25")
	out = addOverlay(out, "assets/logo_wm.png", "-1:48", "(abs(main_h/2)-overlay_w-10):(main_h-overlay_h-10)")

	err := out.Output("pipe:", ffmpeg.KwArgs{"format": "mp4", "movflags": "isml+frag_keyframe", "map": "0:a"},
		ffmpeg.KwArgs{"t": 7}).
		WithOutput(outBuf).
		WithInput(inBuf).
		ErrorToStdOut().
		Run()
	if err != nil {
		log.Fatal(err)
	}
	outFile, err := io.ReadAll(outBuf)
	if err != nil {
		log.Fatal(err)
	}
	return outFile
}

func cropForShort(stream *ffmpeg.Stream) *ffmpeg.Stream {
	w := "abs(ih/2)"
	h := "ih"
	x := "abs(iw/2-ih/2)"
	y := "0"
	return stream.Filter("crop", ffmpeg.Args{fmt.Sprintf("%s:%s:%s:%s", w, h, x, y)})
}

func addOverlay(stream *ffmpeg.Stream, filepath string, scale string, position string) *ffmpeg.Stream {
	overlay := ffmpeg.Input(filepath).Filter("scale", ffmpeg.Args{scale})
	return ffmpeg.Filter([]*ffmpeg.Stream{stream, overlay},
		"overlay", ffmpeg.Args{position})
}
