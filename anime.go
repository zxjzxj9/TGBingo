package main

import (
	"fmt"
	ts "gorgonia.org/tensor"
	"image"
	"io"
)

func animeGAN(reader io.Reader) []byte {
	//encode jpeg to arryy
	img, _, err := image.Decode(reader)
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

	_, ok := imgT.Data().([][][][]float32)
	if !ok {
		fmt.Println("Error conversion image data")
	}

	return nil
}
