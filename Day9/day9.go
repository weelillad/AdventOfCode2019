package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	// Test data
	// testPrograms := [][]int{
	// 	{109,1,204,-1,1001,100,1,100,1008,100,16,101,1006,101,0,99},
	// 	{1102,34915192,34915192,7,4,7,99,0},
	// 	{104,1125899906842624,99},
	// }

	// Actual run
	file, err := os.Open("day9Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()

	program := getProgram(file)

	inputChan := make(chan int)
	defer close(inputChan)
	outputChan := make (chan int)
	go processIntcode(program, inputChan, outputChan)

	// Part 1
	// inputChan <- 1
	//Part 2
	inputChan <- 2

	for output := range outputChan {
		fmt.Print(output, " ")
	}
	fmt.Println("End")
}

// Intcode computer functions

// NOTE: Modifies program array
func processIntcode(program []int, inputChan <-chan int, outputChan chan<- int) int {
	relativeBase := 0
	output := 0

	for position := 0; program[position] != 99; {
		instruction := program[position]
		opcode := instruction % 100
		paramModesArray := getParamModes(instruction)
		switch opcode {
		case 1:
			// fmt.Println("Add instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, operandLeft + operandRight)
			position += 4
		case 2:
			// fmt.Println("Multiply instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, operandLeft * operandRight)
			position += 4
		case 3:
			// fmt.Println("Input instruction")
			program = writeValue(program, program[position+1], getParamModeFromArray(paramModesArray, 0), relativeBase, <-inputChan)
			position += 2
		case 4:
			// fmt.Println("Output instruction")
			output = getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			outputChan <- output
			// fmt.Println("Output: ", output)
			position += 2
		case 5:
			// fmt.Println("Jump-if-true instruction")
			operand := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			if operand != 0 {
				position = getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			} else {
				position += 3
			}
		case 6:
			// fmt.Println("Jump-if-false instruction")
			operand := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			// fmt.Println("Operand: ", operand)
			if operand == 0 {
				position = getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			} else {
				position += 3
			}
		case 7:
			// fmt.Println("Less than instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			if operandLeft < operandRight {
				program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, 1)
			} else {
				program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, 0)
			}
			position += 4
		case 8:
			// fmt.Println("Equals instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			if operandLeft == operandRight {
				program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, 1)
			} else {
				program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, 0)
			}
			position += 4
		case 9:
			// fmt.Println("Relative base update instruction")
			relativeBase += getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			position += 2
		case 99:
			// fmt.Println("End instruction")
			break
		default:
			fmt.Printf("ERROR: Unknown opcode %d", opcode)
		}
		// prettyPrintProgram(program)
	}
	close(outputChan)
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

func getValue(program []int, position, paramMode, relativeBase int) int {
	switch paramMode {
	case 0:
		// position mode
		memAddress := program[position]
		if len(program) > memAddress {
			// fmt.Println("Position mode, memAddress =", memAddress, ", value =", program[memAddress])
			return program[memAddress]
		} else {
			// fmt.Println("Warning: trying to access virgin memory, memAddress =", memAddress)
			return 0
		}
	case 1:
		// immediate mode
			// fmt.Println("Immediate mode, value = ", program[position])
		return program[position]
	case 2:
		// relative mode
		memAddress := program[position]+relativeBase
		if len(program) > memAddress {
			// fmt.Println("Relative mode, memAddress =", memAddress, ", value =", program[memAddress])
			return program[memAddress]
		} else {
			// fmt.Println("Warning: trying to access virgin memory, memAddress =", memAddress)
			return 0
		}
	default:
		log.Fatalf("ERROR: Unknown parameter mode %d", paramMode)
		return -1
	}
}

func writeValue(memory []int, position, paramMode, relativeBase, value int) []int {
	var memAddress int
	switch paramMode {
	case 0:
		memAddress = position
		// fmt.Println("Position mode, memAddress =", memAddress)
	case 1:
		log.Fatal("Writing address cannot be in immediate mode!")
	case 2:
		memAddress = position + relativeBase
		// fmt.Println("Relative mode, memAddress =", memAddress)
	}
	var newMemory []int
	if len(memory) <= memAddress {
		// Expand memory to new position
		// fmt.Println("Expanding memory from size", len(memory), "to", memAddress + 1)
		newMemory = make([]int, memAddress + 1)
		copy(newMemory, memory)
	} else {
		newMemory = memory
	}
	newMemory[memAddress] = value
	return newMemory
}

func getParamModeFromArray(paramModesArray []int, paramPosition int) int {
	if paramPosition < len(paramModesArray) {
		return paramModesArray[paramPosition]
	} else {
		return 0
	}
}
