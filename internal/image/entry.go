package image

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

func Entry() {
	f, err := os.Open("rexpaint_cp437_10x10.png")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	for j := 0; j < 16; j++ {
		for i := 0; i < 16; i++ {
			m := image.NewRGBA(image.Rect(0, 0, 10, 10))
			for y := 0; y < 16; y++ {
				for x := 0; x < 16; x++ {
					c := img.At(i*10+x, j*10+y)
					r, g, b, _ := c.RGBA()
					if r*g*b != 0 {
						m.Set(x, y, c)
					}
				}
			}
			// draw.Draw(m, m.Bounds(), img, image.Point{X: i * 10, Y: j * 10}, draw.Src)

			toImg, err := os.Create(fmt.Sprintf("%v.png", j*16+i))
			if err != nil {
				toImg.Close()
				log.Fatal(err)
			}

			png.Encode(toImg, m)
			toImg.Close()
		}
	}
}
