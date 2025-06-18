package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"strings"
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

func scaleMatrix(matrix [][]int, newWidth, newHeight int) [][]int {
    oldHeight := len(matrix)
    oldWidth := len(matrix[0])
    
    scaled := make([][]int, newHeight)
    for y := 0; y < newHeight; y++ {
        scaled[y] = make([]int, newWidth)
        for x := 0; x < newWidth; x++ {
            // nearest neighbor
            srcY := y * oldHeight / newHeight
            srcX := x * oldWidth / newWidth
            scaled[y][x] = matrix[srcY][srcX]
        }
    }
    return scaled
}

func brightnessToASCII(brightness int) string {
    chars := "@%#*+=-:. "
    
    index := brightness * (len(chars) - 1) / 255
    if index >= len(chars) {
        index = len(chars) - 1
    }
    
    return string(chars[index])
}

func matrixToASCII(matrix [][]int) string {
    var result strings.Builder
    
    for _, row := range matrix {
        for _, brightness := range row {
            result.WriteString(brightnessToASCII(brightness))
        }
        result.WriteString("\n")
    }
    
    return result.String()
}

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Использование: go run main.go <путь_к_изображению> [ширина(опционально)] [высота(опционально)]")
    }
    
    imgPath := os.Args[1]
    
    var asciiWidth, asciiHeight int 

    img, err := openImageFile(imgPath)
    if err != nil {
        log.Fatalf("Ошибка загрузки изображения: %v", err)
    }
    
    bounds := img.Bounds()
    width := bounds.Max.X
    height := bounds.Max.Y

    if width <= 0 || height <= 0 {
        log.Fatal("Некорректные размеры изображения")
    }

    aspectRatio := float64(width) / float64(height)
    fmt.Printf("Пропорции изображения: %.2f (ширина / высота)\n", aspectRatio)

    asciiWidth = 120

    charAspectRatio := 0.5

    asciiHeight = int(float64(asciiWidth) / aspectRatio * charAspectRatio)

    if len(os.Args) >= 3 {
        fmt.Sscanf(os.Args[2], "%d", &asciiWidth)
        asciiHeight = int(float64(asciiWidth) / aspectRatio * charAspectRatio)
    }
    if len(os.Args) >= 4 {
        fmt.Sscanf(os.Args[3], "%d", &asciiHeight)
    }

    if asciiWidth < 1 {
        asciiWidth = 1
    }
    if asciiHeight < 1 {
        asciiHeight = 1
    }
    
    fmt.Printf("Размеры изображения: %d x %d пикселей\n", width, height)
    fmt.Printf("Размеры ASCII: %d x %d символов\n\n", asciiWidth, asciiHeight)
    
    matrix := createPixelMatrix(img)
    
    scaledMatrix := scaleMatrix(matrix, asciiWidth, asciiHeight)
    
    asciiArt := matrixToASCII(scaledMatrix)
    fmt.Print(asciiArt)
    
    // make --save for saving image to txt
    outputFile := "ascii_art.txt"
    err = os.WriteFile(outputFile, []byte(asciiArt), 0644)
    if err != nil {
        log.Printf("Ошибка сохранения файла: %v", err)
    } else {
        fmt.Printf("\nASCII арт сохранен в файл: %s\n", outputFile)
    }
    
}