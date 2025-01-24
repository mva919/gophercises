package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

type storyData struct {
	Title   string
	Story   []string
	Options []struct {
		Text string
		Arc  string
	}
}

func createStoryHTMLTemplate(story storyData) {
	return
}

func main() {
	fileContent, err := os.ReadFile("story.json")
	checkError(err)

	var story map[string]storyData
	jsonUnmarshalErr := json.Unmarshal(fileContent, &story)
	checkError(jsonUnmarshalErr)

	currentArc := "intro"
	stdin := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(story[currentArc].Title, "\n\n")
		for _, paragraph := range story[currentArc].Story {
			fmt.Println("\t", paragraph)
		}

		fmt.Print("\nChoices: \n")
		if len(story[currentArc].Options) == 0 {
			fmt.Println("End of story.")
			return
		}
		for idx, choice := range story[currentArc].Options {
			fmt.Printf("%d - %v: %v\n", idx+1, strings.Join(strings.Split(choice.Arc, "-"), " "), choice.Text)
		}

		var userChoice string
		for {
			fmt.Scan(&userChoice)
			choice, err := strconv.ParseInt(userChoice, 10, 8)
			if err == nil && choice >= 1 && int(choice) <= len(story[currentArc].Options) {
				currentArc = story[currentArc].Options[choice-1].Arc
				break
			}
			stdin.ReadString('\n')
			fmt.Println("Invalid choice pick again.")
		}
	}
}
