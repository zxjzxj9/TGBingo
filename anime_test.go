package main

import (
	"os"
	"testing"
)

func TestAnimeGAN(t *testing.T) {
	reader, err := os.Open("test.png")
	if err != nil {
		t.Error(err)
	}
	outFile := animeGAN(reader)
	writer, err := os.Create("test2.png")
	if err != nil {
		t.Error(err)
	}
	writer.Write(outFile)
	defer reader.Close(), writer.Close()
}
