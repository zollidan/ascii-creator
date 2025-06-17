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

// Функция для масштабирования изображения
func scaleMatrix(matrix [][]int, newWidth, newHeight int) [][]int {
    oldHeight := len(matrix)
    oldWidth := len(matrix[0])
    
    scaled := make([][]int, newHeight)
    for y := 0; y < newHeight; y++ {
        scaled[y] = make([]int, newWidth)
        for x := 0; x < newWidth; x++ {
            // Простое масштабирование nearest neighbor
            srcY := y * oldHeight / newHeight
            srcX := x * oldWidth / newWidth
            scaled[y][x] = matrix[srcY][srcX]
        }
    }
    return scaled
}

// Конвертация яркости в ASCII символ
func brightnessToASCII(brightness int) string {
    // ASCII символы от самого темного к самому светлому
    chars := "@%#*+=-:. "
    
    // Нормализуем яркость к индексу массива символов
    index := brightness * (len(chars) - 1) / 255
    if index >= len(chars) {
        index = len(chars) - 1
    }
    
    return string(chars[index])
}

// Конвертация матрицы пикселей в ASCII
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
    
    // Параметры по умолчанию для ASCII
    asciiWidth := 100
    asciiHeight := 60
    
    // Если указаны размеры, используем их
    if len(os.Args) >= 3 {
        fmt.Sscanf(os.Args[2], "%d", &asciiWidth)
    }
    if len(os.Args) >= 4 {
        fmt.Sscanf(os.Args[3], "%d", &asciiHeight)
    }
    
    img, err := openImageFile(imgPath)
    if err != nil {
        log.Fatalf("Ошибка загрузки изображения: %v", err)
    }
    
    bounds := img.Bounds()
    width := bounds.Max.X
    height := bounds.Max.Y
    
    fmt.Printf("Размеры изображения: %d x %d пикселей\n", width, height)
    fmt.Printf("Размеры ASCII: %d x %d символов\n\n", asciiWidth, asciiHeight)
    
    // Создаем матрицу пикселей
    matrix := createPixelMatrix(img)
    
    // Масштабируем к нужному размеру для ASCII
    scaledMatrix := scaleMatrix(matrix, asciiWidth, asciiHeight)
    
    // Конвертируем в ASCII и выводим
    asciiArt := matrixToASCII(scaledMatrix)
    fmt.Print(asciiArt)
    
    // --save for saving image to txt
    if len(os.Args) >= 5 && os.Args[4] == "--save" {
        outputFile := "ascii_art.txt"
        err := os.WriteFile(outputFile, []byte(asciiArt), 0644)
        if err != nil {
            log.Printf("Ошибка сохранения файла: %v", err)
        } else {
            fmt.Printf("\nASCII арт сохранен в файл: %s\n", outputFile)
        }
    }
}