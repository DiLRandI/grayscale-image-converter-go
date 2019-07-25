package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"runtime"
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

	b := img.Bounds()
	imgSet := image.NewRGBA(b)

	go func() {
		for x := 0; x < b.Max.X; x++ {
			for y := 0; y < b.Max.Y; y++ {
				oldPixel := img.At(x, y)
				r, g, b, _ := oldPixel.RGBA()
				lum := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
				pixel := color.Gray{uint8(lum / 256)}
				imgSet.Set(x, y, pixel)
			}
		}
	}()

	outFile, err := os.Create(fmt.Sprintf("%s/%s", outDir, outFileName))

	if err != nil {
		panic(err)
	}

	defer outFile.Close()
	jpeg.Encode(outFile, imgSet, nil)

	elapsed := time.Since(start)
	log.Printf("Execution time %s", elapsed)
	fmt.Printf("Number of cpus %d \n", runtime.NumCPU())
	fmt.Printf("Number of go rutines %d \n", runtime.NumGoroutine())
}
