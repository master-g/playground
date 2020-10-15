package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
)

var (
	suitList = []string{
		"club",
		"diamond",
		"heart",
		"spade",
	}
	rankList = []string{
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"A",
		"J",
		"K",
		"Q",
		"T",
	}
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: split [INPUT]")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var img image.Image
	img, err = png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	i := 0
	flag := false
	for yi := 0; yi < 6; yi++ {
		for xi := 0; xi < 9; xi++ {
			subImage := img.(SubImager).SubImage(image.Rect(xi*252+2, yi*348+2, (xi+1)*252, (yi+1)*348))
			var outputFile *os.File

			suit := suitList[i%len(suitList)]
			rank := rankList[i/4]

			filename := fmt.Sprintf("%s_%s.png", suit, rank)
			if i == 36 && !flag {
				filename = "background.png"
				flag = true
				i--
			}
			i++

			outputFile, err = os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			err = png.Encode(outputFile, subImage)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
