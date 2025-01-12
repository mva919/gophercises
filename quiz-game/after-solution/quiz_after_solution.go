package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parseProblems(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   line[1],
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}

func main() {
	filename := flag.String("filename", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *filename))
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		exit("Failed to create a new CSV reader")
	}

	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to read the CSV file")
	}

	problems := parseProblems(lines)
	var userAnswer string
	score := 0
	userAnswerCh := make(chan string)

	timer := time.NewTimer(time.Duration(*limit) * time.Second)

problemloop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		go func() {
			fmt.Scanf("%s\n", &userAnswer)
			userAnswerCh <- userAnswer
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break problemloop
		case <-userAnswerCh:
			if userAnswer == p.answer {
				score++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", score, len(problems))
}
