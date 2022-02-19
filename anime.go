package main

import (
	"fmt"
	ts "gorgonia.org/tensor"
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
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// fill the tensor
	imgT := ts.New(ts.Of(ts.Float32), ts.WithShape(1, 3, height, width))

	// set rgb float32 value
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			imgT.SetAt(r/255.0, 0, 0, y, x)
			imgT.SetAt(g/255.0, 0, 1, y, x)
			imgT.SetAt(b/255.0, 0, 2, y, x)
		}
	}

	return nil
}
