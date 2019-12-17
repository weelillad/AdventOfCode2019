package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
)

func main() {
	// Test inputs for Part 1
	// testPrograms := [][]int{
	// 	{3,15,3,16,1002,16,10,16,1,16,15,15,4,15,99,0,0},
	// 	{3,23,3,24,1002,24,10,24,1002,23,-1,23,101,5,23,23,1,24,23,23,4,23,99,0,0},
	// 	{3,31,3,32,1002,32,10,32,1001,31,-2,31,1007,31,0,33,1002,33,7,33,1,33,31,31,1,32,31,31,4,31,99,0,0,0},
	// }
	// testPhaseArrays := [][]int{
	// 	{4,3,2,1,0},
	// 	{0,1,2,3,4},
	// 	{1,0,4,3,2},
	// }
	//
	// // Test input for Part 2
	// testPrograms := [][]int{
	// 	{3, 26, 1001, 26, -4, 26, 3, 27, 1002, 27, 2, 27, 1, 27, 26, 27, 4, 27, 1001, 28, -1, 28, 1005, 28, 6, 99, 0, 0, 5},
	// 	{3, 52, 1001, 52, -5, 52, 3, 53, 1, 52, 56, 54, 1007, 54, 5, 55, 1005, 55, 26, 1001, 54, -5, 54, 1105, 1, 12, 1, 53, 54, 53, 1008, 54, 0, 55, 1001, 55, 1, 55, 2, 53, 55, 53, 4, 53, 1001, 56, -1, 56, 1005, 56, 6, 99, 0, 0, 0, 0, 10},
	// }
	// testPhaseArrays := [][]int{
	// 	{9, 8, 7, 6, 5},
	// 	{9, 7, 8, 5, 6},
	// }

	// Testing scaffold
	// testIndex := 1
	// program := testPrograms[testIndex]
	// phaseCombis := [][]int{testPhaseArrays[testIndex]}

	// Actual run
	file, err := os.Open("day7Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()

	program := getProgram(file)

	// Part 1
	phaseCombis := getPhaseArrayCombinations(5)
	maxOutput := 0
	for _, phaseArray := range phaseCombis {
		output := runAmps(program, phaseArray)
		fmt.Println("Phase Array: ", phaseArray, ", Output: ", output)
		if output > maxOutput {
			maxOutput = output
		}
	}
	fmt.Println("Max output: ", maxOutput)

	// Part 2
	// fmt.Println("Number of orbital transfers", countTransfers(orbitArray))

}

func runAmps(program, phaseArray []int) int {
	ampChannels := make([]chan int, 5)
	for i := 0; i < 5; i++ {
		// Make signal channels
		ampChannels[i] = make(chan int)
		defer close(ampChannels[i])
	}

	// start amps
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		// Make a copy of the program/memory
		memory := make([]int, len(program))
		size := copy(memory, program)
		if size != len(program) {
			log.Fatalf("Failed to copy program")
		}
		var outputChannel chan int
		if i != 4 {
			outputChannel = ampChannels[i+1]
			wg.Add(1)
		} else {
			outputChannel = ampChannels[0]
		}
		go func(phase int, inputChannel <-chan int) {
			defer wg.Done()
			processIntcode(memory, phase, inputChannel, outputChannel)
		}(phaseArray[i], ampChannels[i])
	}
	ampChannels[0] <- 0

	wg.Wait()
	wg.Add(1)
	return <- ampChannels[0]
}

// position is zero-indexed, i.e. last position for length 3 is 2
func getPhaseArrayCombinations(length int) [][]int {
	result := make([][]int, length)
	for i := 0; i < length; i++ {
		result[i] = []int{i}
	}
	if length == 1 {
		return result
	} else {
		candidates := make([]int, length)
		for i := 0; i < length; i++ {
			candidates[i] = i + 5
		}
		return generateCombiRecursive([]int{}, candidates)
	}
}

func generateCombiRecursive(source, candidates []int) [][]int {
	// fmt.Println("generateCombiRecursive: Source ", source, ", candidates ", candidates)
	if len(candidates) == 0 {
		return [][]int{source}
	}
	var result [][]int
	next := make([]int, len(source)+1)
	copy(next, source)
	for i, candidate := range candidates {
		next[len(next)-1] = candidate
		remainingCandidates := make([]int, 0, len(candidates)-1)
		remainingCandidates = append(remainingCandidates, candidates[:i]...)
		// fmt.Println(i, " ", candidates, " ", remainingCandidates)
		if i != len(candidates)-1 {
			remainingCandidates = append(remainingCandidates, candidates[i+1:]...)
		}
		result = append(result, generateCombiRecursive(next, remainingCandidates)...)
	}
	return result
}

// Intcode computer functions

// NOTE: Modifies program array
func processIntcode(program []int, phase int, inputChan <-chan int, outputChan chan<- int) int {
	firstInput := true
	output := 0

	for position := 0; program[position] != 99; {
		instruction := program[position]
		opcode := instruction % 100
		paramModesArray := getParamModes(instruction)
		switch opcode {
		case 1:
			// fmt.Println("Add instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1))
			program[program[position+3]] = operandLeft + operandRight
			position += 4
		case 2:
			// fmt.Println("Multiply instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1))
			program[program[position+3]] = operandLeft * operandRight
			position += 4
		case 3:
			// fmt.Println("Input instruction")
			if firstInput {
				program[program[position+1]] = phase
				firstInput = false
			} else {
				program[program[position+1]] = <-inputChan
			}
			position += 2
		case 4:
			// fmt.Println("Output instruction")
			output = getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			outputChan <- output
			// fmt.Println("Output: ", output)
			position += 2
		case 5:
			// fmt.Println("Jump-if-true instruction")
			operand := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			if operand != 0 {
				position = getValue(program, position+2, getParamModeFromArray(paramModesArray, 1))
			} else {
				position += 3
			}
		case 6:
			// fmt.Println("Jump-if-false instruction")
			operand := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			// fmt.Println("Operand: ", operand)
			if operand == 0 {
				position = getValue(program, position+2, getParamModeFromArray(paramModesArray, 1))
			} else {
				position += 3
			}
		case 7:
			// fmt.Println("Less than instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1))
			if operandLeft < operandRight {
				program[program[position+3]] = 1
			} else {
				program[program[position+3]] = 0
			}
			position += 4
		case 8:
			// fmt.Println("Equals instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0))
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1))
			if operandLeft == operandRight {
				program[program[position+3]] = 1
			} else {
				program[program[position+3]] = 0
			}
			position += 4
		case 99:
			// fmt.Println("End instruction")
			break
		default:
			fmt.Printf("ERROR: Unknown opcode %d", opcode)
		}
		// prettyPrintProgram(program)
	}
	return output
}

func getProgram(file *os.File) []int {
	reader := csv.NewReader(file)
	programString, err := reader.Read()
	if err != nil {
		log.Fatalf("Cannot read input: %s", err)
	}
	program := make([]int, len(programString))
	for i := range programString {
		program[i], err = strconv.Atoi(programString[i])
		if err != nil {
			log.Fatalf("Cannot convert %s to integer: %s", programString[i], err)
		}
	}
	return program
}

func prettyPrintProgram(program []int) {
	for i, val := range program {
		if i == len(program)-1 {
			fmt.Print(val)
		} else {
			fmt.Printf("%d,", val)
		}
	}
	fmt.Print("\n")
}

func getParamModes(instruction int) []int {
	paramModes := instruction / 100
	if paramModes == 0 {
		return []int{}
	}
	paramModesArray := make([]int, int(math.Log10(float64(paramModes))+1))
	for i, index := paramModes, 0; i > 0; index++ {
		paramModesArray[index] = i % 10
		i = i / 10
	}
	return paramModesArray
}

func getValue(program []int, position, paramMode int) int {
	switch paramMode {
	case 0:
		return program[program[position]]
	case 1:
		return program[position]
	default:
		log.Fatalf("ERROR: Unknown parameter mode %d", paramMode)
		return -1
	}
}

func getParamModeFromArray(paramModesArray []int, paramPosition int) int {
	if paramPosition < len(paramModesArray) {
		return paramModesArray[paramPosition]
	} else {
		return 0
	}
}
