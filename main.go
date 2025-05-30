package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

func openImageFile(imageFilename string) (image.Image, error) {
	f, err := os.Open(imageFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func createPixelMatrix(img image.Image) [][]int {
	bounds := img.Bounds()
	width := bounds.Max.X 
	height := bounds.Max.Y

	matrix := make([][]int, height)
	for y := 0; y < height; y++ {
		matrix[y] = make([]int, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			
			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8
			
			brightness := int(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))
			matrix[y][x] = brightness
		}
	}

	return matrix
}

func main() {
	img, err := openImageFile("3d970398314e23a2fa39fca7067e1177.jpg")
	if err != nil {
		log.Fatalf("Error loading image: %v", err)
	}

	// Получаем размеры изображения
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	
	fmt.Printf("Размеры изображения: %d x %d пикселей\n", width, height)
	
	pixelMatrix := createPixelMatrix(img)

}