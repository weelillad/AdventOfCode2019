package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

func main() {
	// Test data
	// testMaps := [][]string{
	// 	{
	// 		// (3,4), 8 detections
	// 		".#..#",
	// 		".....",
	// 		"#####",
	// 		"....#",
	// 		"...##",
	// 	},
	// 	{
	// 		// (5,8), 33 detections
	// 		"......#.#.",
	// 		"#..#.#....",
	// 		"..#######.",
	// 		".#.#.###..",
	// 		".#..#.....",
	// 		"..#....#.#",
	// 		"#..#....#.",
	// 		".##.#..###",
	// 		"##...#..#.",
	// 		".#....####",
	// 	},
	// 	{
	// 		// (1,2), 35 detections
	// 		"#.#...#.#.",
	// 		".###....#.",
	// 		".#....#...",
	// 		"##.#.#.#.#",
	// 		"....#.#.#.",
	// 		".##..###.#",
	// 		"..#...##..",
	// 		"..##....##",
	// 		"......#...",
	// 		".####.###.",
	// 	},
	// 	{
	// 		// (6,3), 41 detections
	// 		".#..#..###",
	// 		"####.###.#",
	// 		"....###.#.",
	// 		"..###.##.#",
	// 		"##.##.#.#.",
	// 		"....###..#",
	// 		"..#.#..#.#",
	// 		"#..#.#.###",
	// 		".##...##.#",
	// 		".....#.#..",
	// 	},
	// 	{
	// 		// (11,13), 210 detections
	// 		".#..##.###...#######",
	// 		"##.############..##.",
	// 		".#.######.########.#",
	// 		".###.#######.####.#.",
	// 		"#####.##.#.##.###.##",
	// 		"..#####..#.#########",
	// 		"####################",
	// 		"#.####....###.#.#.##",
	// 		"##.#################",
	// 		"#####.##.###..####..",
	// 		"..######..##.#######",
	// 		"####.##.####...##..#",
	// 		".#####..#.######.###",
	// 		"##...#.##########...",
	// 		"#.##########.#######",
	// 		".####.#.###.###.#.##",
	// 		"....##.##.###..#####",
	// 		".#.#.###########.###",
	// 		"#.#.#.#####.####.###",
	// 		"###.##.####.##.#..##",
	// 	},
	// }
	// laserPosition := coords{11, 13}

	// Actual run
	file, err := os.Open("day10Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var actualMap []string
	for scanner.Scan() {
		actualMap = append(actualMap, scanner.Text())
	}

	// asteroidList := getAsteroidList(testMaps[4])
	asteroidList := getAsteroidList(actualMap)

	// Part 1
	// var maxVisible int
	// var maxVisibleCoords coords
	// for _, asteroid := range asteroidList {
	// 	offsetList := getOffsetList(asteroidList, asteroid)
	// 	visibleCount := len(getDedupedAngleList(offsetList))
	// 	// fmt.Println("Visible asteroids from", asteroid, ":", visibleCount)
	// 	if visibleCount > maxVisible {
	// 		maxVisible = visibleCount
	// 		maxVisibleCoords = asteroid
	// 	}
	// }
	// fmt.Println("Maximum", maxVisible, "asteroids visible from", maxVisibleCoords)

	// Part 2
	laserPosition := coords{26, 29}
	offsetList := getOffsetList(asteroidList, laserPosition)
	laserOrder := getLaserOrder(offsetList)

	// Debug
	fmt.Println(laserOrder)

	answerOffset := laserOrder[199]
	answerPosition := coords{laserPosition.X + answerOffset.distance.X, laserPosition.Y + answerOffset.distance.Y}
	fmt.Println("Answer:", answerPosition)
}

type coords struct {
	X, Y int
}

func (c coords) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

type coordsArray []coords

func (s coordsArray) Len() int      { return len(s) }
func (s coordsArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func getAsteroidList(asteroidMap []string) []coords {
	var asteroidList []coords
	for y, rowString := range asteroidMap {
		for x, rune := range rowString {
			if rune == '#' {
				asteroidList = append(asteroidList, coords{x, y})
			}
		}
	}

	//Debug
	// fmt.Println(asteroidList)

	return asteroidList
}

type offset struct {
	distance coords
	azimuth float64
}

type offsetArray []offset

func (s offsetArray) Len() int      { return len(s) }
func (s offsetArray) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func getOffsetList(asteroidList []coords, origin coords) []offset {
	var offsetList []offset
	for _, asteroid := range asteroidList {
		if asteroid == origin {
			continue
		}
		distanceX := asteroid.X - origin.X
		distanceY := asteroid.Y - origin.Y
		azimuth := math.Atan2(float64(distanceY), float64(distanceX))
		// Adjust angle such that straight up = smallest value and angles increase clockwise
		if azimuth < -math.Pi / 2 {
			azimuth += 2 * math.Pi
		}
		offsetList = append(
			offsetList,
			offset{
				coords{
					distanceX,
					distanceY,
				},
				azimuth,
			},
		)
	}

	// Sanity check
	if len(asteroidList)-1 != len(offsetList) {
		log.Fatalf("BUG: Distance count does not match asteroid count: Asteroids %v vs Distances %v", len(asteroidList), len(offsetList))
	}

	//Debug
	// fmt.Println(offsetList)

	// sort.Sort(ByDistance{offsetList})

	// //Debug
	// fmt.Println(offsetList)

	return offsetList
}

type ByDistance struct {
	offsetArray
}

func (s ByDistance) Less(i, j int) bool {
	return leftDistanceShorterThanRight(s.offsetArray[i], s.offsetArray[j])
}

func leftDistanceShorterThanRight(left, right offset) bool {
	return (left.distance.X*left.distance.X + left.distance.Y*left.distance.Y) < (right.distance.X*right.distance.X + right.distance.Y*right.distance.Y)
}

type ByAzimuthThenDistance struct {
	offsetArray
}

func (s ByAzimuthThenDistance) Less(i, j int) bool {
	if s.offsetArray[i].azimuth == s.offsetArray[j].azimuth {
		return leftDistanceShorterThanRight(s.offsetArray[i], s.offsetArray[j])
	} else {
		return s.offsetArray[i].azimuth < s.offsetArray[j].azimuth
	}
}

func getDedupedAngleList(offsetList []offset) []offset {
	// Dedupe angles
	var dedupedAngleList []offset
	sort.Sort(ByAzimuthThenDistance{offsetList})
	for i, offset := range offsetList {
		if i == 0 || offset.azimuth != offsetList[i-1].azimuth {
			dedupedAngleList = append(dedupedAngleList, offset)
		}
	}

	// Debug
	// fmt.Println(dedupedAngleList)

	return dedupedAngleList
}

func getLaserOrder(offsetList []offset) []offset {
	var laserOrderList []offset
	var leftoverList []offset

	if len(offsetList) == 0 {
		return []offset{}
	}

	sort.Sort(ByAzimuthThenDistance{offsetList})
	for i, offset := range offsetList {
		if i == 0 || offset.azimuth != offsetList[i-1].azimuth {
			laserOrderList = append(laserOrderList, offset)
		} else {
			leftoverList = append(leftoverList, offset)
		}
	}
	return append(laserOrderList, getLaserOrder(leftoverList)...)
}
