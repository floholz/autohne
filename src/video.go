package src

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
)

type VideoUtils struct {
	silent bool
}

func NewVideoUtils(silent bool) VideoUtils {
	return VideoUtils{silent: silent}
}

type FFmpegStream struct {
	*ffmpeg.Stream
}

func (vu VideoUtils) CreateShort(file []byte) []byte {
	inBuf := bytes.NewBuffer(file)
	outBuf := bytes.NewBuffer(nil)
	out := ffmpeg.Input("pipe:").Filter("scale", ffmpeg.Args{"-1:1080"})

	out = cropForShort(out)
	out = overlayFile(out, "assets/ohnepixel-clipper-mark.png", "500:-1", "25:25", 1.0)
	out = overlayFile(out, "assets/logo_wm.png", "-1:48", "(abs(main_h/2)-overlay_w-10):(main_h-overlay_h-10)", 0.2)

	out = out.Output("pipe:", ffmpeg.KwArgs{"format": "mp4", "movflags": "isml+frag_keyframe", "map": "0:a"},
		ffmpeg.KwArgs{"t": 7}).
		WithOutput(outBuf).
		WithInput(inBuf).
		Silent(vu.silent)

	if !vu.silent {
		out = out.ErrorToStdOut()
	}

	err := out.Run()
	if err != nil {
		log.Fatal(err)
	}
	outFile, err := io.ReadAll(outBuf)
	if err != nil {
		log.Fatal(err)
	}
	return outFile
}

func (vu VideoUtils) CreateShortFromFullVid(file []byte) []byte {
	inBuf := bytes.NewBuffer(file)
	outBuf := bytes.NewBuffer(nil)
	og := ffmpeg.Input("pipe:").Filter("scale", ffmpeg.Args{"-1:1920"})

	split := og.Split()
	split0, split1 := split.Get("0"), split.Get("1")

	split0 = cropForShort(split0).Filter("boxblur", ffmpeg.Args{"50:5"})
	split1 = cropAspectRatio(split1, "1", "1")

	out := overlayStream(split0, split1, "1080:-1", "0:abs(main_h/2-overlay_h/2)")
	out = overlayFile(out, "assets/ohnepixel-clipper-mark.png", "720:-1", "abs(main_w/2-overlay_w/2):64", 1.0)
	out = overlayFile(out, "assets/logo_wm.png", "-1:64", "(abs(main_w)-overlay_w-24):(main_h-overlay_h-24)", 0.2)

	out = out.Output("pipe:", ffmpeg.KwArgs{"format": "mp4", "movflags": "isml+frag_keyframe", "map": "0:a"}).
		WithOutput(outBuf).
		WithInput(inBuf).
		Silent(vu.silent)

	if !vu.silent {
		out = out.ErrorToStdOut()
	}

	err := out.Run()
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
	return stream.Filter("crop", ffmpeg.Args{"1080:ih:abs(iw/2-540):0"})
}

func cropAspectRatio(stream *ffmpeg.Stream, w string, h string) *ffmpeg.Stream {
	return stream.Filter("crop", ffmpeg.Args{
		fmt.Sprintf("abs(iw/%s):abs(ih/%s):abs(iw/%s/2):abs(ih/%s/2)", w, h, w, h),
	})
}

func overlayStream(stream1 *ffmpeg.Stream, stream2 *ffmpeg.Stream, scale string, position string) *ffmpeg.Stream {
	stream2 = stream2.Filter("scale", ffmpeg.Args{scale})
	return ffmpeg.Filter([]*ffmpeg.Stream{stream1, stream2},
		"overlay", ffmpeg.Args{position})
}

func overlayFile(stream *ffmpeg.Stream, filepath string, scale string, position string, opacity float32) *ffmpeg.Stream {
	overlay := ffmpeg.Input(filepath).
		Filter("scale", ffmpeg.Args{scale}).
		ColorChannelMixer(ffmpeg.KwArgs{"aa": fmt.Sprintf("%f", opacity)})
	return ffmpeg.Filter([]*ffmpeg.Stream{stream, overlay},
		"overlay", ffmpeg.Args{position})
}
