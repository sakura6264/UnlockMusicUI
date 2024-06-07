package mmkv

import (
	"encoding/binary"
	"fmt"
	"io"
)

type metadata struct {
	// added in version 0
	crc32 uint32

	// added in version 1
	version  uint32
	sequence uint32 // full write back count

	// added in version 2
	aesVector []byte // random iv for encryption, aes.BlockSize (16 bytes)

	// added in version 3, try to reduce file corruption
	actualSize     uint32
	lastActualSize uint32
	lastCRC32      uint32

	//_reversed []byte // 64 bytes
}

func loadMetadata(rd io.Reader) (*metadata, error) {
	buf := make([]byte, 0x68)
	_, err := io.ReadFull(rd, buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	m := &metadata{}

	m.crc32 = binary.LittleEndian.Uint32(buf[0:4])
	m.version = binary.LittleEndian.Uint32(buf[4:8])

	if m.version >= 1 {
		m.sequence = binary.LittleEndian.Uint32(buf[8:12])
	}

	if m.version >= 2 {
		m.aesVector = buf[12:28]
	}

	if m.version >= 3 {
		m.actualSize = binary.LittleEndian.Uint32(buf[28:32])
		m.lastActualSize = binary.LittleEndian.Uint32(buf[32:36])
		m.lastCRC32 = binary.LittleEndian.Uint32(buf[36:40])
	}

	//m._reversed = buf[40:104]
	return m, nil
}
