package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("day1Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	scanner := bufio.NewScanner(file)

	// sumFuel := part1(scanner)
	sumFuel := part2(scanner)

	log.Printf("Fuel Sum: %d", sumFuel)
}

func part1(scanner *bufio.Scanner) int {
	var sumFuel int
	for scanner.Scan() {
		mass, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalf("Parsing error: %s", err)
		}
		sumFuel += getFuel(mass)
	}
	return sumFuel
}

func part2(scanner *bufio.Scanner) int {
	var sumFuel int
	for scanner.Scan() {
		mass, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatalf("Parsing error: %s", err)
		}
		sumFuel += getFuelRecursive(mass)
	}
	return sumFuel
}

func getFuel(mass int) int {
	return mass/3 - 2
}

func getFuelRecursive(mass int) int {
	massFuel := mass/3 - 2
	if massFuel <= 0 {
		return 0
	} else {
		return massFuel + getFuelRecursive(massFuel)
	}
}
