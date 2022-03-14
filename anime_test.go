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
	defer reader.Close()
}
