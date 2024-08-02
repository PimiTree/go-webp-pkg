package webp

import (
	"log"
	"os"
)

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

		if webp.VP8Chunk.Extended.Animation {
			/*
				+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
				|                      ChunkHeader('ANIM')                      |
				|                                                               |
				+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
				|                       Background Color                        |
				+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
				|          Loop Count           |
				+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
			*/
			webp.getAnimationChunk()
		}
	}
}
