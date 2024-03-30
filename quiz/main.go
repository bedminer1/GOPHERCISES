package main

import ( 
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

type problem struct {
	q string
	a string
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
	}

	// read csv file
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse provided CSV file")
	}

	problems := parseLines(lines)
	score := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = \n", i + 1, p.q)
		var answer string // user's input
		// note that scanf gets rid of spaces so might not be appropriate for string inputs
		fmt.Scanf("%s\n", &answer)
		if answer == p.a {
			score++
			fmt.Println("Correct!")
		} else {
			fmt.Println("Incorrect :(")
		}
	}

	fmt.Printf("Score: %d/%d\n", score, len(problems))
}

// parses csv file into 'problem' struct, flexible with input methods
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]), // edge case: space in answer on csv
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}