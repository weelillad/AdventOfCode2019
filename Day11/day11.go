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
	// Actual run
	file, err := os.Open("day11Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()

	program := getProgram(file)

	inputChan := make(chan int)
	outputChan := make(chan int)
	go func() {
		processIntcode(program, inputChan, outputChan)
		close(inputChan)
	}()

	runRobot(outputChan, inputChan)
}

// Painting robot

type coords struct {
	X, Y int
}

func (c coords) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

type direction int

const (
	DIR_UP direction = iota
	DIR_RIGHT
	DIR_DOWN
	DIR_LEFT
)

type colour int

const (
	COLOUR_BLACK colour = iota
	COLOUR_WHITE
)

func runRobot(inputChannel <-chan int, outputChannel chan<- int) {
	// Create hull map
	mapWidth := 51
	mapHeight := 10
	hullMap := make([][]int, mapHeight)
	for i, _ := range hullMap {
		hullMap[i] = make([]int, mapWidth)
	}
	// Part 1 - create stepMap
	// stepMap := make([][]int, mapSize)
	// for i, _ := range stepMap {
	// 	stepMap[i] = make([]int, mapSize)
	// }

	// initialise robot
	currentPosition := coords{mapWidth / 2 - 23, mapHeight / 2 - 4}
	currentDirection := DIR_UP
	// Part 2 - start on a white square
	hullMap[currentPosition.Y][currentPosition.X] = 1
	// Send camera data to Intcode program
	outputChannel <- hullMap[currentPosition.Y][currentPosition.X]

	var changedPaint, ok bool
	var colourInput int
	paintCount := 0

	// Handle case where Intcode program ends
	defer func() {
		fmt.Println("Robot stopped. Paint count =", paintCount)

		// Debug - Mark position of robot on hullMap
		hullMap[currentPosition.Y][currentPosition.X] = 99

		for _, rowMap := range hullMap {
			for _, val := range rowMap {
				if val == 0 {
					fmt.Print(".")
				} else {
					fmt.Print("X")
				}
			}
			fmt.Print("\n")
		}

		if x := recover(); x != nil {
			fmt.Println("Intcode program ended.")
		}
	}()

	for colourInput, ok = <-inputChannel; ok; colourInput, ok = <-inputChannel {
		rotationInput, ok := <-inputChannel
		if !ok {
			log.Fatal("Received colour but not rotation input")
		}
		// Part 1
		// hullMap, stepMap, currentPosition, currentDirection, changedPaint = takeAction(hullMap, stepMap, currentPosition, currentDirection, colourInput, rotationInput)
		// Part 2
		hullMap, currentPosition, currentDirection, changedPaint = takeAction(hullMap, currentPosition, currentDirection, colourInput, rotationInput)
		if changedPaint {
			paintCount++
		}
		// if paintCount >= 6 {
		if currentPosition.X == 0 || currentPosition.X == mapWidth-1 || currentPosition.Y == 0 || currentPosition.Y == mapHeight-1 {
			fmt.Println("WARNING: Reached border of current hullMap", currentPosition)
			close(outputChannel)
			break
		}

		// Send camera data to Intcode program
		outputChannel <- hullMap[currentPosition.Y][currentPosition.X]
	}
}

// Part 1
// func takeAction(hullMap, stepMap [][]int, currentPosition coords, currentDirection direction, colourInput, rotationInput int) ([][]int, [][]int, coords, direction, bool) {
// Part 2
func takeAction(hullMap [][]int, currentPosition coords, currentDirection direction, colourInput, rotationInput int) ([][]int, coords, direction, bool) {
	// Debug
	// fmt.Println("currentPosition", currentPosition, "currentDirection", currentDirection, "colourInput", colourInput, "rotationInput", rotationInput)

	// Execute paint action
	// Part 1 special case - also mark crossed squares
	// var newSquare bool
	// if stepMap[currentPosition.Y][currentPosition.X] != 1 {
	// 	newSquare = true
	// 	stepMap[currentPosition.Y][currentPosition.X] = 1
	// }
	// Part 2 - normal behavior
	var changedPaint bool
	if hullMap[currentPosition.Y][currentPosition.X] != colourInput {
		changedPaint = true
		hullMap[currentPosition.Y][currentPosition.X] = colourInput
	}

	// Execute rotation action
	var newDirection direction
	switch currentDirection {
	case DIR_UP:
		if rotationInput == 0 {
			newDirection = DIR_LEFT
		} else {
			newDirection = DIR_RIGHT
		}
	case DIR_LEFT:
		if rotationInput == 0 {
			newDirection = DIR_DOWN
		} else {
			newDirection = DIR_UP
		}
	case DIR_DOWN:
		if rotationInput == 0 {
			newDirection = DIR_RIGHT
		} else {
			newDirection = DIR_LEFT
		}
	case DIR_RIGHT:
		if rotationInput == 0 {
			newDirection = DIR_UP
		} else {
			newDirection = DIR_DOWN
		}
	}

	// Execute movement
	var newPosition coords
	switch newDirection {
	case DIR_UP:
		newPosition = coords{currentPosition.X, currentPosition.Y - 1}
	case DIR_RIGHT:
		newPosition = coords{currentPosition.X + 1, currentPosition.Y}
	case DIR_DOWN:
		newPosition = coords{currentPosition.X, currentPosition.Y + 1}
	case DIR_LEFT:
		newPosition = coords{currentPosition.X - 1, currentPosition.Y}
	}

	// Debug
	// fmt.Println("newPosition", newPosition, "newDirection", newDirection, "changedPaint", changedPaint)

	// Part 1
	// return hullMap, stepMap, newPosition, newDirection, newSquare
	// Part 2
	return hullMap, newPosition, newDirection, changedPaint
}

// Intcode computer functions

// NOTE: Modifies program array
func processIntcode(program []int, inputChan <-chan int, outputChan chan<- int) int {
	defer close(outputChan)
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
			program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, operandLeft+operandRight)
			position += 4
		case 2:
			// fmt.Println("Multiply instruction")
			operandLeft := getValue(program, position+1, getParamModeFromArray(paramModesArray, 0), relativeBase)
			operandRight := getValue(program, position+2, getParamModeFromArray(paramModesArray, 1), relativeBase)
			program = writeValue(program, program[position+3], getParamModeFromArray(paramModesArray, 2), relativeBase, operandLeft*operandRight)
			position += 4
		case 3:
			// fmt.Println("Input instruction")
			inputValue := <-inputChan
			program = writeValue(program, program[position+1], getParamModeFromArray(paramModesArray, 0), relativeBase, inputValue)
			// fmt.Println("Input received:", inputValue)
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
	fmt.Println("End instruction encountered. Output = ", output)
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
		memAddress := program[position] + relativeBase
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
		newMemory = make([]int, memAddress+1)
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
