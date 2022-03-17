package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/nfnt/resize"
	"github.com/owulveryck/onnx-go"
	gorgonnx "github.com/owulveryck/onnx-go/backend/x/gorgonnx"
	ts "gorgonia.org/tensor"
	"image"
	"image/color"
	"image/png"
	"io"
)

func animeGAN(reader io.Reader) []byte {
	//encode jpeg to array
	imgRaw, _, err := image.Decode(reader)
	// imgRawHeight := imgRaw.Bounds().Max.Y
	// imgRawWidth := imgRaw.Bounds().Max.X

	// Resize the image to 64x64
	img := resize.Resize(512, 512, imgRaw, resize.Lanczos3)

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

	// Load ONNX model
	// model, err := LoadModel("model/anime.onnx")
	// go:embed:"./face_paint_512_v2_0.onnx"
	var f embed.FS
	data, _ := f.ReadFile("./face_paint_512_v2_0.onnx")
	// fmt.Println(data)
	backend := gorgonnx.NewGraph()
	model := onnx.NewModel(backend)
	err = model.UnmarshalBinary(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = model.SetInput(0, imgT)
	if err != nil {
		fmt.Println(err.Error())
	}

	ret, err := model.GetOutputTensors()
	if err != nil {
		fmt.Println(err.Error())
	}
	imgA := ret[0]
	fmt.Println(imgA)

	imgC := image.NewRGBA(image.Rect(0, 0, 512, 512))

	// set rgb float32 value back to the image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// r, g, b, _ := img.At(x, y).RGBA()
			r, _ := imgA.At(0, 0, y, x)
			g, _ := imgA.At(0, 1, y, x)
			b, _ := imgA.At(0, 2, y, x)
			color := color.RGBA{
				uint8(r.(float32) * 255),
				uint8(g.(float32) * 255),
				uint8(b.(float32) * 255),
				255}
			imgC.Set(x, y, color)
		}
	}
	imgD := resize.Resize(uint(width), uint(height), imgC, resize.Lanczos3)

	buf := bytes.NewBuffer([]byte{})
	err = png.Encode(buf, imgD)
	if err != nil {
		fmt.Println(err.Error())
	}

	return buf.Bytes()
}
