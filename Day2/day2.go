package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("day2Input.txt")
	if err != nil {
		log.Fatalf("Cannot open input file: %s", err)
	}
	defer file.Close()

	program := getProgram(file)

	// output := part1(program)
	// prettyPrintProgram(output)

	fmt.Println(part2(program))

	// Testing of Part 1
	// testProgs := [][]int{
	// 	{1,0,0,0,99},
	// 	{2,3,0,3,99},
	// 	{2,4,4,5,99,0},
	// 	{1,1,1,4,99,5,6,0,99},
	// }
	// for _, testProg := range testProgs {
	// 	prettyPrintProgram(part1(testProg))
	// 	// prettyPrintProgram(part2(testProg))
	// }
}

func part1(program []int) []int {
	for position := 0; program[position] != 99; position += 4{
		currentValue := program[position]
		switch currentValue {
		case 1:
			operandLeft := program[program[position+1]]
			operandRight := program[program[position+2]]
			program[program[position+3]] = operandLeft + operandRight
		case 2:
			operandLeft := program[program[position+1]]
			operandRight := program[program[position+2]]
			program[program[position+3]] = operandLeft * operandRight
		case 99:
			break
		default:
			fmt.Printf("ERROR: Unknown opcode %d", currentValue)
			break
		}
	}
	return program
}

// TODO
func part2(program []int) int {
	var noun, verb int
	output := make([]int, len(program))
	for noun = 12; noun <= 100; noun++ {
		for verb = 2; verb <= 100; verb++ {
			memory := make([]int, len(program))
			copy(memory, program)
			memory[1] = noun
			memory[2] = verb
			output = part1(memory)
			fmt.Printf("Noun: %d, Verb %d, Answer: %d\n", noun, verb, output[0])
			if output[0] == 19690720 {
				fmt.Println("Answer found!!!")
				break
			}
		}
		if output[0] == 19690720 {
			break
		}
	}
	if output[0] == 19690720 {
		return noun * 100 + verb
	} else {
		fmt.Println("Answer not found")
		return 0
	}
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
		if i == len(program) - 1 {
			fmt.Print(val)
		} else {
			fmt.Printf("%d,", val)
		}
	}
	fmt.Print("\n")
}
