package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("day3Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pathSpecs := make([][]pathSpec, 0, 2)
	for scanner.Scan() {
		pathString := scanner.Text()
		pathSpec := parseWirePath(pathString)
		pathSpecs = append(pathSpecs, pathSpec)
	}

	// Test sets
	// pathSpecs = append(pathSpecs, parseWirePath("R8,U5,L5,D3"))
	// pathSpecs = append(pathSpecs, parseWirePath("U7,R6,D4,L4"))

	// pathSpecs = append(pathSpecs, parseWirePath("R75,D30,R83,U83,L12,D49,R71,U7,L72"))
	// pathSpecs = append(pathSpecs, parseWirePath("U62,R66,U55,R34,D71,R55,D58,R83"))

	// pathSpecs = append(pathSpecs, parseWirePath("R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51"))
	// pathSpecs = append(pathSpecs, parseWirePath("U98,R91,D20,R16,D67,R40,U7,R15,U6,R7"))

	wireMap, zeroVal := plotWirePathTravelCost(pathSpecs[0])
	answer := overlayWirePathTravelCost(pathSpecs[1], wireMap, zeroVal)

	fmt.Printf("Answer: %v\n", answer)
}

type coord struct {
	X int
	Y int
}

func getManhattanDistance(location coord, zeroVal int) int {
	return int(math.Abs(float64(location.X - zeroVal)) + math.Abs(float64(location.Y - zeroVal)))
}

type pathSpec struct {
	Direction string
	Distance  int
}

func parseWirePath(pathString string) []pathSpec {
	wirePath := []pathSpec{}
	pathStringFragments := strings.Split(pathString, ",")
	for _, fragment := range pathStringFragments {
		distance, err := strconv.Atoi(fragment[1:])
		if err != nil {
			log.Fatalf("Cannot convert %s to integer", fragment[1:])
		}
		wirePath = append(wirePath, pathSpec{
			Direction: fragment[0:1],
			Distance:  distance,
		})
	}
	return wirePath
}

func plotWirePathBool(wirePath []pathSpec) ([][]bool, int) {
	zeroVal := 11001
	wireMap := make([][]bool, 22001)
	for i := range wireMap {
		wireMap[i] = make([]bool, 22001)
	}
	pointer := coord{zeroVal, zeroVal}
	for _, spec := range wirePath {
		switch spec.Direction {
		case "U":
			finalY := pointer.Y + spec.Distance
			for i := pointer.Y + 1; i <= finalY; i++ {
				wireMap[pointer.X][i] = true
			}
			pointer.Y = finalY
		case "D":
			finalY := pointer.Y - spec.Distance
			for i := pointer.Y - 1; i >= finalY; i-- {
				wireMap[pointer.X][i] = true
			}
			pointer.Y = finalY
		case "L":
			finalX := pointer.X - spec.Distance
			for i := pointer.X - 1; i >= finalX; i-- {
				wireMap[i][pointer.Y] = true
			}
			pointer.X = finalX
		case "R":
			finalX := pointer.X + spec.Distance
			for i := pointer.X + 1; i <= finalX; i++ {
				wireMap[i][pointer.Y] = true
			}
			pointer.X = finalX
		default:
			log.Fatalf("Invalid direction: %s", spec.Direction)
		}
	}
	return wireMap, zeroVal
}

func plotWirePathTravelCost(wirePath []pathSpec) ([][]int, int) {
	zeroVal := 11001
	wireMap := make([][]int, 22001)
	for i := range wireMap {
		wireMap[i] = make([]int, 22001)
	}
	pointer := coord{zeroVal, zeroVal}
	distance := 0
	for _, spec := range wirePath {
		switch spec.Direction {
		case "U":
			finalY := pointer.Y + spec.Distance
			for i := pointer.Y + 1; i <= finalY; i++ {
				distance++
				if wireMap[pointer.X][i] == 0 {
					wireMap[pointer.X][i] = distance
				}
			}
			pointer.Y = finalY
		case "D":
			finalY := pointer.Y - spec.Distance
			for i := pointer.Y - 1; i >= finalY; i-- {
				distance++
				if wireMap[pointer.X][i] == 0{
					wireMap[pointer.X][i] = distance
				}
			}
			pointer.Y = finalY
		case "L":
			finalX := pointer.X - spec.Distance
			for i := pointer.X - 1; i >= finalX; i-- {
				distance++
				if wireMap[i][pointer.Y] == 0{
					wireMap[i][pointer.Y] = distance
				}
			}
			pointer.X = finalX
		case "R":
			finalX := pointer.X + spec.Distance
			for i := pointer.X + 1; i <= finalX; i++ {
				distance++
				if wireMap[i][pointer.Y] == 0 {
					wireMap[i][pointer.Y] = distance
				}
			}
			pointer.X = finalX
		default:
			log.Fatalf("Invalid direction: %s", spec.Direction)
		}
	}
	return wireMap, zeroVal
}

func overlayWirePath(wirePath []pathSpec, wireMap [][]bool, zeroVal int) int {
	shortestDistance := 999999
	pointer := coord{zeroVal, zeroVal}
	for _, spec := range wirePath {
		switch spec.Direction {
		case "U":
			finalY := pointer.Y + spec.Distance
			for i := pointer.Y + 1; i <= finalY; i++ {
				if wireMap[pointer.X][i] == true {
					intersectionPoint := coord{pointer.X, i}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestManhattanDistance(shortestDistance, intersectionPoint, zeroVal)
				}
			}
			pointer.Y = finalY
		case "D":
			finalY := pointer.Y - spec.Distance
			for i := pointer.Y - 1; i >= finalY; i-- {
				if wireMap[pointer.X][i] == true {
					intersectionPoint := coord{pointer.X, i}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestManhattanDistance(shortestDistance, intersectionPoint, zeroVal)
				}
			}
			pointer.Y = finalY
		case "L":
			finalX := pointer.X - spec.Distance
			for i := pointer.X - 1; i >= finalX; i-- {
				if wireMap[i][pointer.Y] == true {
					intersectionPoint := coord{i, pointer.Y}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestManhattanDistance(shortestDistance, intersectionPoint, zeroVal)
				}
			}
			pointer.X = finalX
		case "R":
			finalX := pointer.X + spec.Distance
			for i := pointer.X + 1; i <= finalX; i++ {
				if wireMap[i][pointer.Y] == true {
					intersectionPoint := coord{i, pointer.Y}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestManhattanDistance(shortestDistance, intersectionPoint, zeroVal)
				}
			}
			pointer.X = finalX
		default:
			log.Fatalf("Invalid direction: %s", spec.Direction)
		}
	}
	return shortestDistance
}

func overlayWirePathTravelCost(wirePath []pathSpec, wireMap [][]int, zeroVal int) int {
	shortestDistance := 999999
	travelDistance := 0
	pointer := coord{zeroVal, zeroVal}
	for _, spec := range wirePath {
		switch spec.Direction {
		case "U":
			finalY := pointer.Y + spec.Distance
			for i := pointer.Y + 1; i <= finalY; i++ {
				travelDistance++
				if wireMap[pointer.X][i] > 0 {
					intersectionPoint := coord{pointer.X, i}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestSignalDelay(shortestDistance, travelDistance + wireMap[pointer.X][i])
				}
			}
			pointer.Y = finalY
		case "D":
			finalY := pointer.Y - spec.Distance
			for i := pointer.Y - 1; i >= finalY; i-- {
				travelDistance++
				if wireMap[pointer.X][i] > 0 {
					intersectionPoint := coord{pointer.X, i}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestSignalDelay(shortestDistance, travelDistance + wireMap[pointer.X][i])
				}
			}
			pointer.Y = finalY
		case "L":
			finalX := pointer.X - spec.Distance
			for i := pointer.X - 1; i >= finalX; i-- {
				travelDistance++
				if wireMap[i][pointer.Y] > 0 {
					intersectionPoint := coord{i, pointer.Y}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestSignalDelay(shortestDistance, travelDistance + wireMap[i][pointer.Y])
				}
			}
			pointer.X = finalX
		case "R":
			finalX := pointer.X + spec.Distance
			for i := pointer.X + 1; i <= finalX; i++ {
				travelDistance++
				if wireMap[i][pointer.Y] > 0 {
					intersectionPoint := coord{i, pointer.Y}
					fmt.Printf("Intersection found at %v\n", intersectionPoint)
					shortestDistance = getShortestSignalDelay(shortestDistance, travelDistance + wireMap[i][pointer.Y])
				}
			}
			pointer.X = finalX
		default:
			log.Fatalf("Invalid direction: %s", spec.Direction)
		}
	}
	return shortestDistance
}

func getShortestManhattanDistance(shortestDistance int, location coord, zeroVal int) int {
	intersectionDistance := getManhattanDistance(location, zeroVal)
	fmt.Printf("Intersection found at %v, distance %d\n", location, intersectionDistance)
	if (intersectionDistance < shortestDistance) {
		return intersectionDistance
	} else {
		return shortestDistance
	}
}

func getShortestSignalDelay(shortestDistance int, newDistance int) int {
	if (newDistance < shortestDistance) {
		return newDistance
	} else {
		return shortestDistance
	}
}
