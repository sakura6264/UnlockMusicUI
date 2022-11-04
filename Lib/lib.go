package main
import "C"
import (
	"os"
	"path/filepath"
	"strings"
	//"fmt"
	"github.com/unlock-music/cli/algo/common"
	_ "github.com/unlock-music/cli/algo/kgm"
	_ "github.com/unlock-music/cli/algo/kwm"
	_ "github.com/unlock-music/cli/algo/ncm"
	_ "github.com/unlock-music/cli/algo/qmc"
	_ "github.com/unlock-music/cli/algo/tm"
	_ "github.com/unlock-music/cli/algo/xm"
)
//export DecFile
func DecFile(inFile *C.char, outDir *C.char, SkipNoop C.int) C.int {
	// 0:no error 1:no decoder -1:decode error -2:output error
	inputFile := C.GoString(inFile)
	outputDir := C.GoString(outDir)
	//fmt.Println("inputFile:",inputFile)
	//fmt.Println("outputDir:",outputDir)
	allDec := common.GetDecoder(inputFile, SkipNoop!=0)
	if len(allDec) == 0 {
		return 1
	}
	//fmt.Println("allDecLen:",len(allDec))
	file, err := os.ReadFile(inputFile)
	if err != nil {
		return -1
	}
	var dec common.Decoder
	for _, decFunc := range allDec {
		dec = decFunc(file)
		if err := dec.Validate(); err == nil {
			break
		}
		dec = nil
	}
	if dec == nil {
		return -1
	}
	if err := dec.Decode(); err != nil {
		return -1
	}
	outData := dec.GetAudioData()
	outExt := dec.GetAudioExt()
	//fmt.Println("outExt:",outExt)
	if outExt == "" {
		if ext, ok := common.SniffAll(outData); ok {
			outExt = ext
		} else {
			outExt = ".mp3"
		}
	}
	filenameOnly := strings.TrimSuffix(filepath.Base(inputFile), filepath.Ext(inputFile))
	//fmt.Println("filenameOnly:",filenameOnly)
	outPath := filepath.Join(outputDir, filenameOnly+outExt)
	//fmt.Println("outPath:",outPath)
	err = os.WriteFile(outPath, outData, 0644)
	if err != nil {
		return -2
	}
	return 0
}