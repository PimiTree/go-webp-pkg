package webp

import "fmt"

// Info prints file elements to stdout. Mostly usage to debug pkg and check the image structure consistency
func (webp *WEBP) Info() {
	fmt.Printf("Path: %s \n Header: \n	%s \n	%d \n	%s \nVP8Chunk: \n	Header: %s \n	ChunkSize: %d \n",
		webp.FilePath,
		webp.Header.RIFFChunk.Value,
		webp.Header.FileSize.Value+4,
		webp.Header.WEBPChunk.Value,
		webp.VP8Chunk.Header.Value,
		webp.VP8Chunk.Header.ChunkSize,
	)

	if webp.VP8Chunk.Header.Value == VP8Extended {
		fmt.Printf("Extended fields: \n	ILEXAR: %08b \n	ICC: %t \n	Alpha: %t \n	Exif: %t \n	XMP: %t \n	Animation: %t \n",
			webp.VP8Chunk.Extended.ILEXAR,
			webp.VP8Chunk.Extended.ICC,
			webp.VP8Chunk.Extended.Alpha,
			webp.VP8Chunk.Extended.Exif,
			webp.VP8Chunk.Extended.XMP,
			webp.VP8Chunk.Extended.Animation,
		)
	}
	if webp.VP8Chunk.Extended.Animation {
		fmt.Printf("AnimationChunk: \n	Header: %s \n	ChunkSize: %d \n	BackgroundColor: %v \n 	LoopCount: %d \n		LastPosition: %d \n",
			webp.Animation.Header.Value,
			webp.Animation.Header.ChunkSize,
			webp.Animation.BackgroundColor,
			webp.Animation.LoopCount,
			webp.Animation.LastPosition,
		)
		for _, frame := range webp.Animation.Frames {
			fmt.Printf("	AnimationFrame: \n		Header: %s \n		ChunkSize: %d \n		FrameX: %d \n		FrameY: %d \n		FrameWidth: %d \n		FrameHeight:  %d \n		FrameDuration:   %d \n		BlendingMethod: %t \n		DisposalMethod: %t \n",
				frame.Header.Value,
				frame.Header.ChunkSize,
				frame.FrameX,
				frame.FrameY,
				frame.FrameWidth,
				frame.FrameHeight,
				frame.FrameDuration,
				frame.BlendingMethod,
				frame.DisposalMethod,
			)
		}
	}
}
