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
	file, err := os.Open("day5Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()

	program := getProgram(file)

	// Part 1
	// processIntcode(program, 1)

	// Part 2
	processIntcode(program, 5)

	// Testing for Part 1
	// testProgs := [][]int{
	// 	{3, 0, 4, 0, 99},
	// 	{1002,4,3,4,33},
	// 	{1101,100,-1,4,0},
	// }

	// Testing for Part 2
	// testProgs := [][]int{
		// {3,9,8,9,10,9,4,9,99,-1,8},
		// {3,9,7,9,10,9,4,9,99,-1,8},
		// {3,3,1108,-1,8,3,4,3,99},
		// {3,3,1107,-1,8,3,4,3,99},
		// {3,12,6,12,15,1,13,14,13,4,13,99,-1,0,1,9},
		// {3,3,1105,-1,9,1101,0,0,12,4,12,99,1},
		// {3,21,1008,21,8,20,1005,20,22,107,8,21,20,1006,20,31,1106,0,36,98,0,0,1002,21,125,20,4,20,1105,1,46,104,999,1105,1,46,1101,1000,1,20,4,20,1105,1,46,98,99},
	// }

	// for _, testProg := range testProgs {
		// processIntcode(testProg, 8)
		// processIntcode(testProg, 0)
		// processIntcode(testProg, 10)
	// }
}

func processIntcode(program []int, input int) []int {
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
			program[program[position+1]] = input
			position += 2
		case 4:
			// fmt.Println("Output instruction")
			fmt.Println("Output: ", getValue(program, position+1, getParamModeFromArray(paramModesArray, 0)))
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
	return program
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
