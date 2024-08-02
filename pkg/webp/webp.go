package webp

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	VP8Lossy    = "VP8 "
	VP8LossLess = "VP8L "
	VP8Extended = "VP8X"
)

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
		Header struct {
			Raw       []byte
			Value     string
			ChunkSize uint32
		}
		Extended struct {
			ILEXAR    string
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
}

func (webp *WEBP) GetData() {
	data, err := os.ReadFile(webp.FilePath)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}
	webp.Data = data
	// https://developers.google.com/speed/webp/docs/riff_container
	// WEBP file header 12 bytes represents:
	/*
		 0                   1                   2                   3
		 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|      'R'      |      'I'      |      'F'      |      'F'      |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                           File Size                           |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|      'W'      |      'E'      |      'B'      |      'P'      |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	*/
	webp.getWEBPHeader()
	webp.getVP8ChunkHeader()

	if webp.VP8Chunk.Header.Value == VP8Extended {
		///*
		//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//	|                                                               |
		//	|                   WebP file header (12 bytes)                 |
		//	|                                                               |
		//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//	|                      ChunkHeader('VP8X')                      |
		//	|                                                               |
		//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//	|Rsv|I|L|E|X|A|R|                   Reserved                    |
		//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//	|          Canvas Width Minus One               |             ...
		//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//	...  Canvas Height Minus One    |
		//	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		//*/

		webp.getVP8ExtendedChunk()
		webp.getVP8ExtendedDimensions()
	}

}

// getWEBPHeader takes the whole WEBP header
func (webp *WEBP) getWEBPHeader() {
	webp.getRiffChunk()
	webp.getFileSizeChunk()
	webp.getWEBPChunk()
}

// getRiffChunk takes RIFF: 4 bytes(uint8, byte):  The ASCII characters 'R', 'I', 'F', 'F'.
func (webp *WEBP) getRiffChunk() {
	webp.Header.RIFFChunk.Raw = webp.Data[0:4]
	webp.Header.RIFFChunk.Value = string(webp.Header.RIFFChunk.Raw)
}

// getFileSizeChunk takes File Size: 4 bytes (uint32): The size of the file in bytes, starting at offset 8
func (webp *WEBP) getFileSizeChunk() {
	webp.Header.FileSize.Raw = webp.Data[4:8]
	/*
		The file size in the header
		is the total size of the chunks that follow plus 4 bytes for the 'WEBP' FourCC.
		The file SHOULD NOT contain any data after the data specified by File Size.
	*/
	webp.Header.FileSize.Value = binary.LittleEndian.Uint32(webp.Header.FileSize.Raw) + 4
}

// getWEBPChunk takes WEBP 4 bytes (uint8, byte):The ASCII characters 'W', 'E', 'B', 'P'
func (webp *WEBP) getWEBPChunk() {
	webp.Header.WEBPChunk.Raw = webp.Data[8:12]
	webp.Header.WEBPChunk.Value = string(webp.Header.WEBPChunk.Raw)
}

// getVP8ChunkHeader takes the VP8 ChunkHeader it can be 'VP8 ', 'VP8L', 'VP8X'
func (webp *WEBP) getVP8ChunkHeader() {
	webp.VP8Chunk.Header.Raw = webp.Data[12:20]
	webp.VP8Chunk.Header.Value = string(webp.VP8Chunk.Header.Raw[0:4])
	webp.VP8Chunk.Header.ChunkSize = binary.LittleEndian.Uint32(webp.VP8Chunk.Header.Raw[4:8])
}

// getVP8ExtendedChunk takes the extended information акщь VP8X ChunkHeader
func (webp *WEBP) getVP8ExtendedChunk() {
	webp.VP8Chunk.Extended.ILEXAR = fmt.Sprintf("%08b", webp.Data[20])

	m := make(map[string]bool, 2)
	m["0"] = false
	m["1"] = true

	webp.VP8Chunk.Extended.ICC = m[string(webp.VP8Chunk.Extended.ILEXAR[2])]
	webp.VP8Chunk.Extended.Alpha = m[string(webp.VP8Chunk.Extended.ILEXAR[3])]
	webp.VP8Chunk.Extended.Exif = m[string(webp.VP8Chunk.Extended.ILEXAR[4])]
	webp.VP8Chunk.Extended.XMP = m[string(webp.VP8Chunk.Extended.ILEXAR[5])]
	webp.VP8Chunk.Extended.Animation = m[string(webp.VP8Chunk.Extended.ILEXAR[6])]
}

// getVP8ExtendedDimensions takes canvas Width and Height
func (webp *WEBP) getVP8ExtendedDimensions() {
	CanvasWidthByteSlice := webp.Data[24:27]
	CanvasHeightByteSlice := webp.Data[27:30]

	webp.Canvas.Width = uint32(CanvasWidthByteSlice[2])<<16 | uint32(CanvasWidthByteSlice[1])<<8 | uint32(CanvasWidthByteSlice[0]) + 1
	webp.Canvas.Height = uint32(CanvasHeightByteSlice[2])<<16 | uint32(CanvasHeightByteSlice[1])<<8 | uint32(CanvasHeightByteSlice[0]) + 1
}

// Info prints file elements to stdout. Mostly usage to debug pkg and check the image structure consistency
func (webp *WEBP) Info() {
	fmt.Printf("Path: %s \n Header: \n	%s \n	%d \n	%s \nVP8Chunk: \n	Header: %s \n	ChunkSize: %d \n",
		webp.FilePath,
		webp.Header.RIFFChunk.Value,
		webp.Header.FileSize.Value,
		webp.Header.WEBPChunk.Value,
		webp.VP8Chunk.Header.Value,
		webp.VP8Chunk.Header.ChunkSize,
	)

	if webp.VP8Chunk.Header.Value == VP8Extended {
		fmt.Printf("Extended fields: \n	ILEXAR: %s \n	ICC: %t \n	Alpha: %t \n	Exif: %t \n	XMP: %t \n	Animation: %t \n",
			webp.VP8Chunk.Extended.ILEXAR,
			webp.VP8Chunk.Extended.ICC,
			webp.VP8Chunk.Extended.Alpha,
			webp.VP8Chunk.Extended.Exif,
			webp.VP8Chunk.Extended.XMP,
			webp.VP8Chunk.Extended.Animation,
		)
	}
}
