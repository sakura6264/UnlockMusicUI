package meta

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bodgit/sevenzip"
)

//go:embed ffmpeg.7z
var ffmpegZip []byte

func writeTempExec(pattern string, binary []byte) (string, error) {
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer f.Close()
	_, err = f.Write(binary)
	if err != nil {
		return "", fmt.Errorf("fail to write executable: %v", err)
	}
	if err := f.Chmod(os.ModePerm); err != nil {
		return "", fmt.Errorf("fail to chmod: %v", err)
	}
	return f.Name(), nil
}

func un7zip() (ffmpeg []byte, ffprobe []byte, err error) {
	r, err := sevenzip.NewReader(bytes.NewReader(ffmpegZip), int64(len(ffmpegZip)))
	if err != nil {
		return
	}
	extract := func(f *sevenzip.File) error {
		if f.FileInfo().IsDir() {
			return nil
		}
		isFFmpeg := strings.HasSuffix(f.Name, "ffmpeg.exe")
		isFFprobe := strings.HasSuffix(f.Name, "ffprobe.exe")
		if !isFFmpeg && !isFFprobe {
			return nil
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		buffer, err := io.ReadAll(rc)
		if err != nil {
			return err
		}
		if isFFmpeg {
			ffmpeg = buffer
		} else if isFFprobe {
			ffprobe = buffer
		}

		return nil
	}
	for idx := range r.File {
		if len(ffmpeg) > 0 && len(ffprobe) > 0 {
			break
		}
		err = extract(r.File[idx])
		if err != nil {
			return
		}
	}
	return
}

var (
	ffmpegPath  string
	ffprobePath string
)

func FFmpegPath() string { return ffmpegPath }

func FFprobePath() string { return ffprobePath }

var inited bool = false

func Init() {
	if inited {
		return
	}
	var err error
	ffmpeg, ffprobe, err := un7zip()
	if err != nil {
		panic(fmt.Errorf("failed to unarchive ffmpeg: %v", err))
	}
	ffmpegPath, err = writeTempExec("ffmpeg_temp.exe", ffmpeg)
	if err != nil {
		panic(fmt.Errorf("failed to write ffmpeg_temp: %v", err))
	}
	ffprobePath, err = writeTempExec("ffprobe_temp.exe", ffprobe)
	if err != nil {
		panic(fmt.Errorf("failed to write ffprobe_temp: %v", err))
	}
	inited = true
}
