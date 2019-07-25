package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"time"
)

type ImageSet interface {
	Set(x, y int, c color.Color)
}

func main() {
	start := time.Now()

	if len(os.Args) != 4 {
		panic("Expected 3 arguments,\n 1. Path to image. \n 2. Path to output folder. \n 3. Path to output file name.")
	}

	pathToFile := os.Args[1]
	outDir := os.Args[2]
	outFileName := os.Args[3]
	file, err := os.Open(pathToFile)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)

	if err != nil {
		panic(err)
	}

	outFile, err := os.Create(fmt.Sprintf("%s/%s", outDir, outFileName))

	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	b := img.Bounds()
	imgSet := processImage(0, b.Max.X, 0, b.Max.Y, img)
	jpeg.Encode(outFile, imgSet, nil)

	elapsed := time.Since(start)
	log.Printf("Execution time %s", elapsed)
}

func processImage(startX, maxX, startY, maxY int, img image.Image) *image.RGBA {
	imgSet := image.NewRGBA(img.Bounds())

	for x := startX; x < maxX; x++ {
		for y := startY; y < maxY; y++ {
			oldPixel := img.At(x, y)
			r, g, b, _ := oldPixel.RGBA()
			lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			pixel := color.Gray{uint8(lum / 256)}
			imgSet.Set(x, y, pixel)
		}
	}

	return imgSet
}
