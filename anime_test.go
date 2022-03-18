package main

import (
	"embed"
	"os"
	"testing"
)

func TestAnimeGAN(t *testing.T) {
	// go:embed:"./face_paint_512_v2_0.onnx"
	var f embed.FS
	reader, err := f.Open("./test.png")

	// reader, err := os.Open("test.png")
	if err != nil {
		t.Error(err)
	}
	outFile := animeGAN(reader)
	writer, err := os.Create("test2.png")
	if err != nil {
		t.Error(err)
	}
	writer.Write(outFile)
	defer reader.Close()
	defer writer.Close()
}
