package main

import (
	"fmt"
	"image"
	"io"
)

func animeGAN(r io.Reader) []byte {
	//encode jpeg to arryy
	img, _, err := image.Decode(r)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(img)
	return nil
}
