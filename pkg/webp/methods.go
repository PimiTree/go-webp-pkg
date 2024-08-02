package webp

import (
	"encoding/binary"
	"imageConverterFromRaw/pkg/dataConverter"
)

// getWEBPHeader takes the whole WEBP header.
func (webp *WEBP) getWEBPHeader() {
	webp.getRiffChunk()
	webp.getFileSizeChunk()
	webp.getWEBPChunk()
}

// getRiffChunk takes RIFF: 4 bytes(uint8, byte):  The ASCII characters 'R', 'I', 'F', 'F'. Raw data slice [0:4]
func (webp *WEBP) getRiffChunk() {
	webp.Header.RIFFChunk.Raw = webp.Data[0:4]
	webp.Header.RIFFChunk.Value = string(webp.Header.RIFFChunk.Raw)
}

// getFileSizeChunk takes File Size: 4 bytes (uint32): The size of the file in bytes, starting at offset 8.  Raw data slice [4:8]
func (webp *WEBP) getFileSizeChunk() {
	webp.Header.FileSize.Raw = webp.Data[4:8]
	/*
		The file size in the header
		is the total size of the chunks that follow plus 4 bytes for the 'WEBP' FourCC.
		The file SHOULD NOT contain any data after the data specified by File Size.
	*/
	webp.Header.FileSize.Value = binary.LittleEndian.Uint32(webp.Header.FileSize.Raw) + 4
}

// getWEBPChunk takes WEBP 4 bytes (uint8, byte):The ASCII characters 'W', 'E', 'B', 'P'. Raw data slice [8:12]
func (webp *WEBP) getWEBPChunk() {
	webp.Header.WEBPChunk.Raw = webp.Data[8:12]
	webp.Header.WEBPChunk.Value = string(webp.Header.WEBPChunk.Raw)
}

// getVP8ChunkHeader takes the VP8 ChunkHeader it can be 'VP8 ', 'VP8L', 'VP8X'. Raw data slice [12:20]
func (webp *WEBP) getVP8ChunkHeader() {
	webp.VP8Chunk.Header = chunkHeaderCreate(webp.Data[12:20])
}

// getVP8ExtendedChunk takes the extended information акщь VP8X ChunkHeader.
func (webp *WEBP) getVP8ExtendedChunk() {
	ILEXAR := webp.Data[20]
	webp.VP8Chunk.Extended.ILEXAR = ILEXAR

	webp.VP8Chunk.Extended.ICC = isBitwiseTrue(ILEXAR, 5)
	webp.VP8Chunk.Extended.Alpha = isBitwiseTrue(ILEXAR, 4)
	webp.VP8Chunk.Extended.Exif = isBitwiseTrue(ILEXAR, 3)
	webp.VP8Chunk.Extended.XMP = isBitwiseTrue(ILEXAR, 2)
	webp.VP8Chunk.Extended.Animation = isBitwiseTrue(ILEXAR, 1)
}

// getVP8ExtendedDimensions takes canvas Width and Height. Raw data slice [24:30]
func (webp *WEBP) getVP8ExtendedDimensions() {
	CanvasWidthByteSlice := webp.Data[24:27]
	CanvasHeightByteSlice := webp.Data[27:30]

	webp.Canvas.Width = dataConverter.FromUint24TRoUint32(CanvasWidthByteSlice) + 1
	webp.Canvas.Height = dataConverter.FromUint24TRoUint32(CanvasHeightByteSlice) + 1
}

// getAnimationChunk takes Animation chunk header and its fields for "ANIM". Raw data slice [30:41]
func (webp *WEBP) getAnimationChunk() {
	webp.getANIM()

	webp.getANMF(44)

}
func (webp *WEBP) getANIM() {
	webp.Animation.Header = chunkHeaderCreate(webp.Data[30:38])

	// BackgroundColor of ChunkHeader('ANIM') has [Blue, Green, Red, Alpha] byte order
	webp.Animation.BackgroundColor[0] = webp.Data[38]                       // R
	webp.Animation.BackgroundColor[0] = webp.Data[39]                       // G
	webp.Animation.BackgroundColor[0] = webp.Data[40]                       // B
	webp.Animation.BackgroundColor[0] = webp.Data[41]                       // Alpha
	webp.Animation.LoopCount = binary.LittleEndian.Uint16(webp.Data[42:44]) // 0 == infinity
}

// getAnimationChunk takes Animation chunk header and its fields for "ANMF". Raw data [44:68]
func (webp *WEBP) getANMF(position uint32) {
	reserved := webp.Data[position+22]

	webp.Animation.Frames = append(webp.Animation.Frames, AnimationFrame{
		Header:         chunkHeaderCreate(webp.Data[position : position+8]),
		FrameX:         dataConverter.FromUint24TRoUint32(webp.Data[position+8 : position+11]),    // The X coordinate of the upper left corner of the frame is Frame X * 2.
		FrameY:         dataConverter.FromUint24TRoUint32(webp.Data[position+11 : position+14]),   // The Y coordinate of the upper left corner of the frame is Frame Y * 2.
		FrameWidth:     dataConverter.FromUint24TRoUint32(webp.Data[position+14:position+17]) + 1, // data holds as Frame Width Minus One
		FrameHeight:    dataConverter.FromUint24TRoUint32(webp.Data[position+17:position+20]) + 1, // data holds as Frame Width Minus One
		FrameDuration:  dataConverter.FromUint24TRoUint32(webp.Data[position+20 : position+23]),
		Reserved:       reserved,
		BlendingMethod: isBitwiseTrue(reserved, 0),
		DisposalMethod: isBitwiseTrue(reserved, 1),
		FrameData:      FrameData{},
	})

	nextPosition := position + 8 + webp.Animation.Frames[len(webp.Animation.Frames)-1].Header.ChunkSize
	webp.Animation.LastPosition = nextPosition
	if nextPosition >= webp.Header.FileSize.Value {
		return
	}

	if string(webp.Data[nextPosition:nextPosition+4]) == "ANMF" {
		webp.getANMF(nextPosition)
	}

}

func isBitwiseTrue(b byte, s int) bool {
	return ((b >> s) & 1) == 1
}
