package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

func main() {
	// Test data
	// image := []int{0,2,2,2,1,1,2,2,2,2,1,2,0,0,0,0}
	// const imgWidth, imgHeight int = 2, 2

	// Actual run
	const imgWidth, imgHeight int = 25, 6
	input, err := ioutil.ReadFile("day8Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	image := convertImageStringToIntArray(string(input[:len(input)-1]))
	imageLayers := splitImageIntoLayers(image, imgWidth, imgHeight)

	// Part 1
	// for i, layer := range imageLayers {
	// 	count0, count1, count2 := countDigits(layer)
	// 	fmt.Println("Layer ", i, ": ", count0, " 0s, ", count1, " 1s, ", count2, " 2s")
	// }

	// Part 2
	fullImage := getFullImage(imageLayers, imgWidth * imgHeight)
	printImage(fullImage, imgWidth)
}

func convertImageStringToIntArray(image string) []int {
	intArray := make([]int, len(image))
	for i := range image {
		integer, err := strconv.Atoi(string(image[i]))
		if err != nil {
			log.Fatalf("Failed to convert image bit: %v", image[i])
		}
		intArray[i] = integer
	}
	return intArray
}

func splitImageIntoLayers(image []int, width, height int) [][]int {
	imageSize := len(image)
	layerSize := width * height
	if imageSize%layerSize != 0 {
		log.Fatalf("Image cannot be cleanly split into layers! Image size: %d, layer size: %d", imageSize, layerSize)
	}

	pointer := 0
	layers := make([][]int, 0, imageSize/layerSize)
	for pointer < imageSize {
		layer := make([]int, layerSize)
		copy(layer, image[pointer:pointer+layerSize])
		layers = append(layers, layer)
		pointer += layerSize
	}

	return layers
}

func countDigits(image []int) (int, int, int) {
	var count0, count1, count2 int
	for _, pixel := range image {
		switch pixel {
		case 0:
			count0++
		case 1:
			count1++
		case 2:
			count2++
		}
	}

	return count0, count1, count2
}

func getFullImage(imageLayers [][]int, layerSize int) []int {
	fullImage := make([]int, layerSize)

	for i := 0; i < layerSize; i++ {
		layerPointer := 0
		for imageLayers[layerPointer][i] == 2 {
			layerPointer++
		}
		fullImage[i] = imageLayers[layerPointer][i]
	}

	return fullImage
}

func printImage(image []int, imgWidth int) {
	for i, pixel := range image {
		fmt.Print(pixel)
		if (i + 1) % imgWidth == 0 {
			fmt.Print("\n")
		}
	}
}
