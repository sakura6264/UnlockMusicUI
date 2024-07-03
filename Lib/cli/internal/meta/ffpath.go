package meta

import (
	"os/exec"
	"path/filepath"
	"syscall"
)

var (
	ffmpegPath  string
	ffprobePath string
)

func FFmpegPath() string { return ffmpegPath }

func FFprobePath() string { return ffprobePath }

var inited bool = false

func ForceInit(exedir string) {
	if inited {
		return
	}

	// check ffmpeg in the same directory
	ffmpeg := filepath.Join(exedir, "ffmpeg.exe")
	ffprobe := filepath.Join(exedir, "ffprobe.exe")
	commandffmpeg := exec.Command(ffmpeg, "-version")
	commandffprobe := exec.Command(ffprobe, "-version")
	commandffmpeg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	commandffprobe.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	errffmpeg := commandffmpeg.Run()
	errffprobe := commandffprobe.Run()
	if errffmpeg == nil && errffprobe == nil {
		ffmpegPath = ffmpeg
		ffprobePath = ffprobe
		inited = true
		return
	}
	// check system ffmpeg
	commandffmpeg = exec.Command("ffmpeg.exe", "-version")
	commandffprobe = exec.Command("ffprobe.exe", "-version")
	commandffmpeg.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	commandffprobe.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	errffmpeg = commandffmpeg.Run()
	errffprobe = commandffprobe.Run()
	if errffmpeg == nil && errffprobe == nil {
		ffmpegPath = "ffmpeg.exe"
		ffprobePath = "ffprobe.exe"
		inited = true
		return
	}
}
