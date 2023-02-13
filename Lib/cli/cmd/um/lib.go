package main
import "C"
import (
	"os"
	"path/filepath"
	"bytes"
	"io"
	"context"
	"time"
	"errors"
	"strings"
	//"fmt"
	"unlock-music.dev/cli/algo/common"
	_ "unlock-music.dev/cli/algo/kgm"
	_ "unlock-music.dev/cli/algo/kwm"
	_ "unlock-music.dev/cli/algo/ncm"
	_ "unlock-music.dev/cli/algo/qmc"
	_ "unlock-music.dev/cli/algo/tm"
	_ "unlock-music.dev/cli/algo/xiami"
	_ "unlock-music.dev/cli/algo/ximalaya"
	"unlock-music.dev/cli/internal/ffmpeg"
	"unlock-music.dev/cli/internal/sniff"
	"unlock-music.dev/cli/internal/utils"
)
type processor struct {
	outputDir string

	skipNoopDecoder bool
	removeSource    bool
	updateMetadata  bool
	overwriteOutput bool
}

//export DecFile
func DecFile(inFile *C.char, outDir *C.char, SkipNoop C.int) C.int {
	// 0:no error 1:no decoder -1:decode error -2:output error
	inputFile := C.GoString(inFile)
	outputDir := C.GoString(outDir)
	//fmt.Println("inputFile:",inputFile)
	//fmt.Println("outputDir:",outputDir)
	allDec := common.GetDecoder(inputFile, SkipNoop!=0)
	p := processor{
		outputDir:     outputDir,
		updateMetadata: true,
		removeSource:  false,
		overwriteOutput: false,
		skipNoopDecoder: SkipNoop!=0,
	}
	if len(allDec) == 0 {
		return 1
	}
	//fmt.Println("allDecLen:",len(allDec))
	file, err := os.Open(inputFile)
	
	if err != nil {
		return -1
	}
	defer file.Close()
	decParams := &common.DecoderParams{
		Reader:    file,
		Extension: filepath.Ext(inputFile),
		FilePath:  inputFile,
	}
	var dec common.Decoder
	for _, decFunc := range allDec {
		dec = decFunc(decParams)
		if err := dec.Validate(); err == nil {
			break
		}
		dec = nil
	}
	if dec == nil {
		return -1
	}
	params := &ffmpeg.UpdateMetadataParams{}

	header := bytes.NewBuffer(nil)
	_, err = io.CopyN(header, dec, 64)
	if err != nil {
		return -1
	}
	audio := io.MultiReader(header, dec)
	params.AudioExt = sniff.AudioExtensionWithFallback(header.Bytes(), ".mp3")
	if p.updateMetadata {
		if audioMetaGetter, ok := dec.(common.AudioMetaGetter); ok {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// since ffmpeg doesn't support multiple input streams,
			// we need to write the audio to a temp file.
			// since qmc decoder doesn't support seeking & relying on ffmpeg probe, we need to read the whole file.
			// TODO: support seeking or using pipe for qmc decoder.
			params.Audio, err = utils.WriteTempFile(audio, params.AudioExt)
			if err != nil {
				return -1
			}
			defer os.Remove(params.Audio)

			params.Meta, _ = audioMetaGetter.GetAudioMeta(ctx)

			if params.Meta == nil { // reset audio meta if failed
				audio, err = os.Open(params.Audio)
				if err != nil {
					return -1
				}
			}
		}
	}
	if p.updateMetadata && params.Meta != nil {
		if coverGetter, ok := dec.(common.CoverImageGetter); ok {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			cover, err := coverGetter.GetCoverImage(ctx) 
			imgExt, ok := sniff.ImageExtension(cover)
			if err == nil && ok {
				params.AlbumArtExt = imgExt
				params.AlbumArt = cover
			}
		}
	}
	inFilename := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	outPath := filepath.Join(p.outputDir, inFilename+params.AudioExt)
	_, e := os.Stat(outPath)

	if e == nil {
		return -2
	} else if !errors.Is(e, os.ErrNotExist) {
		return -2
	}
	if params.Meta == nil {
		outFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return -2
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, audio); err != nil {
			return -2
		}
		outFile.Close()

	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		if err := ffmpeg.UpdateMeta(ctx, outPath, params); err != nil {
			return -1
		}
	}
	return 0
}