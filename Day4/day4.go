package main

import (
	"fmt"
	"log"
	"strconv"
)

func main() {
	// fmt.Println(validPassword(111111))
	// fmt.Println(validPassword(223450))
	// fmt.Println(validPassword(123789))

	// fmt.Println(validPassword(112233))
	// fmt.Println(validPassword(123444))
	// fmt.Println(validPassword(111122))

	answer := 0
	for candidate := 246515; candidate <= 739105; candidate++ {
		if validPassword(candidate) {
			fmt.Println("Valid password", candidate)
			answer++
		}
	}
	fmt.Println("Number of valid passwords", answer)
}

func hasTwoOrMoreSameAdjacentDigits(candidate int) bool {
	digits := strconv.Itoa(candidate)
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			return true
		}
	}
	return false
}

func hasExactlyTwoSameAdjacentDigits(candidate int) bool {
	digits := strconv.Itoa(candidate)
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			// Find the end of the sequence
			seqLength := 2
			j := i+2
			for ; j < len(digits) && digits[j] == digits[i]; j++ {
				seqLength++
			}
			if seqLength == 2 {
				fmt.Println("Exactly 2 adjacent digits found: ", digits[i:i+1])
				return true
			} else {
				// Move the pointer to the end of the sequence
				i = j - 1
			}
		}
	}
	return false
}

func hasIncreasingDigits(candidate int) bool {
	digits := strconv.Itoa(candidate)
	digitArray := make([]int, len(digits))
	var err error
	for i := 0; i < len(digits); i++ {
		digitArray[i], err = strconv.Atoi(digits[i : i+1])
		if err != nil {
			log.Fatalf("Cannot convert %v to integer", digits[i])
		}
		if i == 0 {
			continue
		}
		if digitArray[i] < digitArray[i-1] {
			return false
		}
	}
	return true
}

func validPassword(candidate int) bool {
	return hasIncreasingDigits(candidate) && hasExactlyTwoSameAdjacentDigits(candidate)
}
