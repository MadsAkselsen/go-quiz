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
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
		os.Exit(1)
	}
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLine(lines)

	// NewTimer writes to a channel, 'timer.C', after the set amount of time
	// time.Duration converts timeLimit to the type that is compatible with Newtimer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0

problemLoop:
	for i, p := range problems {
		// show question
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerCh := make(chan string)

		// Using go routing for user input because we're also waiting
		// for the timer, and we don't want them to block each other
		go func() {
			var answer string
			// get user answer from terminal
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTimes up! You scored %d out of %d.\n", correct, len(problems))
			break problemLoop
		case answer := <- answerCh:
			// check if answer is correct
			if answer == p.a {
				correct++
			}

			if (i == len(problems)-1) {
				fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			}
		}
	}
}

func parseLine(lines [][]string) []problem {
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

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}