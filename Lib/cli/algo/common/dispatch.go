package common

import (
	"io"
	"path/filepath"
	"strings"

	"go.uber.org/zap"
)

type DecoderParams struct {
	Reader    io.ReadSeeker // required
	Extension string        // required, source extension, eg. ".mp3"

	FilePath string // optional, source file path

	Logger *zap.Logger // required
}
type NewDecoderFunc func(p *DecoderParams) Decoder

type decoderItem struct {
	noop    bool
	decoder NewDecoderFunc
}

var DecoderRegistry = make(map[string][]decoderItem)

func RegisterDecoder(ext string, noop bool, dispatchFunc NewDecoderFunc) {
	DecoderRegistry[ext] = append(DecoderRegistry[ext],
		decoderItem{noop: noop, decoder: dispatchFunc})
}

func GetDecoder(filename string, skipNoop bool) (rs []NewDecoderFunc) {
	ext := strings.ToLower(strings.TrimLeft(filepath.Ext(filename), "."))
	for _, dec := range DecoderRegistry[ext] {
		if skipNoop && dec.noop {
			continue
		}
		rs = append(rs, dec.decoder)
	}
	return
}
