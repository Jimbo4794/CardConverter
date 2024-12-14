package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nfnt/resize"
)

/*
	How it works:
	- Takes an image and converts to gray-scale. We do this so we only have to work with a
	single digit representation of a pixel.
	- Next convert to a binary image by using a pixels value as density value.
	- We fill the cluster (in this case 9x9) randomly based of the density value. (Monte
	Carlo-ing here is a little lazy as you could probably make this more crisp by calculating
	the best way to "Spread" the pixels across the new cluster)
	- We then resize the image to avoid huge image sizes with all the new pixels. 1pixel -> 9x9 cluster
	- Ive set this to 380 based of the printer limitations

	@author: github.com/Jimbo4794
*/

var clusterSize = 9      // increasing the cluster size will produce smoother images but take longer
var concurrentLimit = 10 // The threading limit.

func main() {
	filePaths := getFilePaths("./input")
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrentLimit)

	for _, i := range filePaths {
		wg.Add(1)
		go convertImage(i, &wg, semaphore)
	}
	wg.Wait()

	log.Printf("Finished! Converted %v cards", len(filePaths))
}

func convertImage(imagePath string, wg *sync.WaitGroup, semaphore chan struct{}) {
	log.Printf("To convert image %s...", path.Base(imagePath))
	defer wg.Done()
	semaphore <- struct{}{}
	img := openImage(imagePath)
	grayImg := toGrayscale(img)
	binaryImg := generateBinaryImage(grayImg, clusterSize)
	saveImage(strings.Replace(imagePath, "input", "output", -1), *binaryImg)
	<-semaphore
}

func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			grayValue := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			grayImg.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}
	return grayImg
}

func grayToCluster(grayValue uint8, clusterSize int) [][]uint8 {
	cluster := make([][]uint8, clusterSize)
	for i := range cluster {
		cluster[i] = make([]uint8, clusterSize)
	}

	totalPixels := clusterSize * clusterSize
	blackPixels := int(float64(totalPixels) * (1.0 - float64(grayValue)/255.0))

	for blackPixels > 0 {
		x, y := rand.Intn(clusterSize), rand.Intn(clusterSize)
		if cluster[y][x] == 0 {
			cluster[y][x] = 1
			blackPixels--
		}
	}
	return cluster
}

func generateBinaryImage(grayImg *image.Gray, clusterSize int) *image.Image {
	bounds := grayImg.Bounds()
	newWidth := (bounds.Max.X - bounds.Min.X) * clusterSize
	newHeight := (bounds.Max.Y - bounds.Min.Y) * clusterSize

	binaryImg := image.NewGray(image.Rect(0, 0, newWidth, newHeight))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayValue := grayImg.GrayAt(x, y).Y
			cluster := grayToCluster(grayValue, clusterSize)
			for i := 0; i < clusterSize; i++ {
				for j := 0; j < clusterSize; j++ {
					newX := x*clusterSize + j
					newY := y*clusterSize + i
					if cluster[i][j] == 1 {
						binaryImg.SetGray(newX, newY, color.Gray{Y: 0})
					} else {
						binaryImg.SetGray(newX, newY, color.Gray{Y: 255})
					}
				}
			}
		}
	}

	width := uint(380) // bounded by the printer spec
	height := uint(532)

	finalImage := resize.Resize(width, height, binaryImg, resize.Lanczos2)

	return &finalImage
}

func getFilePaths(filePath string) []string {
	files := make([]string, 0)
	entries, err := os.ReadDir(filePath)
	if err != nil {
		panic(err)
	}

	for _, e := range entries {
		if e.IsDir() {
			files = append(files, getFilePaths(path.Join(filePath, e.Name()))...)
		} else {
			files = append(files, path.Join(filePath, e.Name()))
		}
	}

	return files
}

func openImage(imagePath string) image.Image {
	inFile, _ := os.Open(imagePath)
	defer inFile.Close()

	img, err := jpeg.Decode(inFile)
	if err != nil {
		panic(err)
	}
	return img
}

func saveImage(filePath string, img image.Image) {
	log.Printf("%s converting. Saving to %s", path.Base(filePath), filePath)

	dir := filepath.Dir(filePath)
	os.MkdirAll(dir, os.ModePerm)

	outfile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	err = jpeg.Encode(outfile, img, nil)
	if err != nil {
		panic(err)
	}
}
