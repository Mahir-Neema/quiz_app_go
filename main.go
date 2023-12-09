package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file containing questions and answers")

	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")

	flag.Parse()
	file, err := os.Open(*csvFilename)

	if err != nil {
		fmt.Printf("Failed to open csv file: %s\n", *csvFilename)
		os.Exit(1)
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()

	if err != nil {
		fmt.Println("Failed to parse the provided csv file")
	}

	// fmt.Println(lines)
	problems := parseLines(lines)
	// fmt.Println(problems)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))

	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}

	return ret
}

type problem struct {
	q string
	a string
}
