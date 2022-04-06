package main

import (
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", http.HandlerFunc(DefaultHandler)))
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	draw(w)
}

func draw(w io.Writer) {
	data := []int{100, 200, 260, 180}
	width, height := len(data)*80+20, 300
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)

	//set background white
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			img.Set(i, j, color.RGBA{255, 255, 255, 255})
		}
	}

	for i, d := range data {
		for x := i*80 + 20; x < (i+1)*80; x++ {
			for y := height; y > (height - d); y-- {
				img.Set(x, y, color.RGBA{180, 180, 250, 255})
			}
		}
	}

	png.Encode(w, img)
}
