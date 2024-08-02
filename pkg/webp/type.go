package webp

import "encoding/binary"

type WEBP struct {
	FilePath string
	Data     []byte
	Header   struct {
		RIFFChunk struct {
			Raw   []byte
			Value string
		}
		FileSize struct {
			Raw   []byte
			Value uint32
		}
		WEBPChunk struct {
			Raw   []byte
			Value string
		}
	}
	VP8Chunk struct {
		Header   ChunkHeader
		Extended struct {
			ILEXAR    byte
			ICC       bool
			Alpha     bool
			Exif      bool
			XMP       bool
			Animation bool
		}
	}
	Canvas struct {
		Width     uint32
		Height    uint32
		Bitstream []byte
	}
	Animation struct {
		Header          ChunkHeader
		BackgroundColor [4]byte // [Blue, Green, Red, Alpha] byte order
		LoopCount       uint16
		Frames          []AnimationFrame
		LastPosition    uint32
	}
}

type AnimationFrame struct {
	Header         ChunkHeader
	FrameX         uint32 // uint24
	FrameY         uint32 // uint24
	FrameWidth     uint32 // uint24
	FrameHeight    uint32 // uint24
	FrameDuration  uint32 // uint24
	Reserved       byte
	BlendingMethod bool
	DisposalMethod bool
	FrameData      FrameData
}

type FrameData struct {
}

type ChunkHeader struct {
	Raw       []byte
	Value     string
	ChunkSize uint32
}

func chunkHeaderCreate(raw []byte) ChunkHeader {
	return ChunkHeader{
		Raw:       raw,
		Value:     string(raw[0:4]),
		ChunkSize: binary.LittleEndian.Uint32(raw[4:8]),
	}
}
