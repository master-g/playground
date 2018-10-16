package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
)

const (
	WIDTH  = 405
	HEIGHT = 585
)

func Entry() {
	f, err := os.Open("unoCards.png")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	for j := 0; j < 6; j++ {
		for i := 0; i < 10; i++ {
			m := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
			draw.Draw(m, m.Bounds(), img, image.Point{X: 2 + i*WIDTH + 3*i, Y: j * HEIGHT}, draw.Src)

			toImg, err := os.Create(fmt.Sprintf("%v-%v.png", j, i))
			if err != nil {
				toImg.Close()
				log.Fatal(err)
			}

			png.Encode(toImg, m)
			toImg.Close()
		}
	}
}
