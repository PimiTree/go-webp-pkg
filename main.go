package main

import "imageConverterFromRaw/pkg/webp"

func main() {
	Test1 := webp.WEBP{}
	Test1.FilePath = "./assets/WebpAnimatedAlpha.webp"
	Test1.GetData()
	Test1.Info()
}
