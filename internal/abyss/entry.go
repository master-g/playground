package abyss

import (
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

const (
	SCREEN_WIDTH  = 1024
	SCREEN_HEIGHT = 576
)

// Entry for this package
func Entry() {
	// load images
	img1 := gocv.IMRead("loading.png", gocv.IMReadGrayScale)
	img2 := gocv.IMRead("screen.png", gocv.IMReadGrayScale)

	// init SIFT detector
	sift := contrib.NewSIFT()

	// find key points and descriptors with SIFT
	kp1, des1 := sift.DetectAndCompute(img1, gocv.NewMat())
	kp2, des2 := sift.DetectAndCompute(img2, gocv.NewMat())

	FLANN_INDEX_KDTREE := 1
	flann :=

	window := gocv.NewWindow("Abyss")
	window.ResizeWindow(SCREEN_WIDTH, SCREEN_HEIGHT)
	window.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowNormal)

	for {
		window.IMShow(img2)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
