package mmkv

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_loadMetadata(t *testing.T) {
	file, err := os.Open("./testdata/mmkv.default.crc")
	require.NoError(t, err)

	meta, err := loadMetadata(file)
	require.NoError(t, err)

	assert.Equal(t, uint32(3), meta.version)
	assert.Equal(t, uint32(1), meta.sequence)

	assert.Equal(t, uint32(28), meta.actualSize)
	assert.Equal(t, uint32(197326043), meta.crc32)

	assert.Equal(t, uint32(4), meta.lastActualSize)
	assert.Equal(t, uint32(1285129681), meta.lastCRC32)

	assert.Equal(t, bytes.Repeat([]byte{0x00}, 16), meta.aesVector)
}
