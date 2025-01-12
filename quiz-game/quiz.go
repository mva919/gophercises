package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	csv_filename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	limit := flag.Uint("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csv_filename)
	check(err)
	defer file.Close()

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	check(err)

	var user_answer string
	var score uint16
	timer_started := false

	for idx, line := range data {
		fmt.Printf("Problem #%d: %s = ", idx+1, line[0])
		fmt.Scan(&user_answer)
		if !timer_started {
			timer_started = true
			go func() {
				timer := time.NewTimer(time.Duration(*limit) * time.Second)
				<-timer.C
				fmt.Printf("\nYou scored %d out of %d.\n", score, len(data))
				os.Exit(0)
			}()
		}
		if user_answer == line[1] {
			score++
		}
	}
	fmt.Printf("You scored %d out of %d.\n", score, len(data))
}
