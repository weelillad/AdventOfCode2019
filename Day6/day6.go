package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// Test input for Part 1
	// testOrbits := []string{
	// 	"COM)B",
	// 	"B)C",
	// 	"C)D",
	// 	"D)E",
	// 	"E)F",
	// 	"B)G",
	// 	"G)H",
	// 	"D)I",
	// 	"E)J",
	// 	"J)K",
	// 	"K)L",
	// }
	//
	// // Test input for Part 2
	// testOrbits := []string{
	// 	"COM)B",
	// 	"B)C",
	// 	"C)D",
	// 	"D)E",
	// 	"E)F",
	// 	"B)G",
	// 	"G)H",
	// 	"D)I",
	// 	"E)J",
	// 	"J)K",
	// 	"K)L",
	// 	"K)YOU",
	// 	"I)SAN",
	// }

	// orbitArray := make([]orbit, 0, len(testOrbits))
	// for _, orbitString := range testOrbits {
	// 	orbitArray = append(orbitArray, parseOrbit(orbitString))
	// }

	file, err := os.Open("day6Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	orbitArray := []orbit{}
	for scanner.Scan() {
		orbitArray = append(orbitArray, parseOrbit(scanner.Text()))
	}

	// Part 1
	// fmt.Println("Orbit Count", countOrbits(orbitArray))

	// Part 2
	fmt.Println("Number of orbital transfers", countTransfers(orbitArray))
}

type orbit struct {
	orbitee string
	orbiter string
}

func parseOrbit(orbitString string) orbit {
	orbitArray := strings.Split(orbitString, ")")
	if len(orbitArray) != 2 {
		log.Fatalf("Orbit parsing error: %s", orbitString)
	}
	orbitObj := orbit{
		orbitee: orbitArray[0],
		orbiter: orbitArray[1],
	}

	// Debug
	// fmt.Println(orbitObj.orbitee, ")", orbitObj.orbiter)

	return orbitObj
}

func countOrbits(inputArray []orbit) int {
	// Copy inputArray before making changes
	orbitArray := inputArray

	orbitCount := 0
	for len(orbitArray) > 0 {
		// Debug
		// fmt.Println("Number of orbits", len(orbitArray))

		orbitCount += len(orbitArray)
		orbitArray = pruneOrbitRoots(orbitArray)
	}

	return orbitCount
}

func pruneOrbitRoots(inputArray []orbit) []orbit {
	// Get map of orbiters, for O(1) lookup
	orbiterMap := make(map[string]bool, len(inputArray))
	for _, orbit := range inputArray {
		orbiterMap[orbit.orbiter] = true
	}

	// Append orbit that are also orbiters i.e. not orbit roots
	outputArray := make([]orbit, 0, len(inputArray))
	for _, orbit := range inputArray {
		_, ok := orbiterMap[orbit.orbitee]
		if ok {
			outputArray = append(outputArray, orbit)
		} else {
			// Debug
			// fmt.Println("Pruned orbit", orbit)
		}
	}

	return outputArray
}

func countTransfers(inputArray []orbit) int {
	// Copy inputArray before making changes
	orbitArray := inputArray

	orbitArrayLength := 0
	for len(orbitArray) != orbitArrayLength {
		orbitArrayLength = len(orbitArray)
		orbitArray = pruneOrbitLeaves(orbitArray)
	}

	// Hack the truck until the common joint is reached
	orbitArray = pruneOrbitTrunk(orbitArray)

	// Debug
	fmt.Println(orbitArray)
	// Sort orbitArray from YOU to SAN
	orbitObjMap := make(map[string][]string, len(orbitArray) + 2)
	for _, orbit := range orbitArray {
		orbiteeLinks, ok := orbitObjMap[orbit.orbitee]
		if !ok {
			orbitObjMap[orbit.orbitee] = []string{orbit.orbiter}
		} else {
			orbitObjMap[orbit.orbitee] = append(orbiteeLinks, orbit.orbiter)
		}

		orbiterLinks, ok := orbitObjMap[orbit.orbiter]
		if !ok {
			orbitObjMap[orbit.orbiter] = []string{orbit.orbitee}
		} else {
			orbitObjMap[orbit.orbiter] = append(orbiterLinks, orbit.orbitee)
		}
	}
	pointer := "YOU"
	beforePointer := ""
	for pointer != "SAN" {
		fmt.Printf("%s > ", pointer)
		linkedNodes := orbitObjMap[pointer]
		switch len(linkedNodes) {
		case 1:
			beforePointer = pointer
			pointer = linkedNodes[0]
		case 2:
			if beforePointer == linkedNodes[0] {
				if beforePointer == linkedNodes[1] {
					log.Fatalf("Something went wrong...")
				}
				beforePointer = pointer
				pointer = linkedNodes[1]
			} else {
				beforePointer = pointer
				pointer = linkedNodes[0]
			}
		default:
			log.Fatalf("Something went wrong...")
		}
	}
	fmt.Printf("%s\n", pointer)

	// Return the number of orbits, minus 2 (YOU and SAN's orbit)
	return len(orbitArray) - 2
}

func pruneOrbitLeaves(inputArray []orbit) []orbit {
	// Get map of orbitees, for O(1) lookup
	orbiteeMap := make(map[string]bool, len(inputArray))
	for _, orbit := range inputArray {
		orbiteeMap[orbit.orbitee] = true
	}

	// Append orbit that are also orbitees i.e. not orbit leaves
	// Exception: SAN and YOU
	outputArray := make([]orbit, 0, len(inputArray))
	for _, orbit := range inputArray {
		_, ok := orbiteeMap[orbit.orbiter]
		if ok || orbit.orbiter == "SAN" || orbit.orbiter == "YOU" {
			outputArray = append(outputArray, orbit)
		} else {
			// Debug
			fmt.Println("Pruned orbit", orbit)
		}
	}

	return outputArray
}

// Orbit trunk is defined as the orbit chain of objects, starting from the root object, that have only 1 orbiter
func pruneOrbitTrunk(inputArray []orbit) []orbit {
	// Get map of orbiters, for O(1) lookup
	orbiterMap := make(map[string]bool, len(inputArray))
	// Get map of orbitees and how many direct orbiters each has
	orbiteeToNumOrbitersMap := make(map[string]int, len(inputArray))
	for _, orbit := range inputArray {
		orbiterMap[orbit.orbiter] = true
		orbiteeToNumOrbitersMap[orbit.orbitee] += 1
	}

	// Copy inputArray before making changes
	outputArray := inputArray
	tempArray := make([]orbit, 0, len(outputArray))
	for true {
		// Drop roots with only 1 orbiter
		for _, orbit := range outputArray {
			if orbiterMap[orbit.orbitee] || orbiteeToNumOrbitersMap[orbit.orbitee] > 1 {
				tempArray = append(tempArray, orbit)
			} else {
				// Debug
				// fmt.Println("Pruned orbit", orbit)
				// Update orbiter map
				orbiterMap[orbit.orbiter] = false
			}
		}

		if len(outputArray) == len(tempArray) {
			// Nothing was removed in this round, we're done
			return tempArray
		} else {
			outputArray = tempArray
			tempArray = make([]orbit, 0, len(outputArray))
		}
	}

	return nil
}
