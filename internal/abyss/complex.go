package abyss

import (
	"image"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
)

func oldCode() {
	var err error
	var imgData []byte
	imgData, err = screenCapture()
	if err != nil {
		log.Fatalf("failed to capture screen, %v", err)
	}
	log.Infof("screen captured, %v bytes for png", len(imgData))

	var img gocv.Mat
	img, err = gocv.IMDecode(imgData, gocv.IMReadColor)
	if err != nil {
		log.Fatalf("failed to decode image data, %v", err)
	}

	window := gocv.NewWindow("Abyss")
	window.ResizeWindow(SCREEN_WIDTH, SCREEN_HEIGHT)
	window.SetWindowProperty(gocv.WindowPropertyAutosize, gocv.WindowNormal)

	gocv.Resize(img, &img, image.Point{X: SCREEN_WIDTH, Y: SCREEN_HEIGHT}, 0, 0, gocv.InterpolationLinear)

	for {
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func screenCapture() (pngData []byte, err error) {
	cmd := exec.Command("adb", "exec-out", "screencap", "-p")
	pngData, err = cmd.CombinedOutput()
	return
}
