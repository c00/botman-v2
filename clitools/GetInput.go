package clitools

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func SetInput(label string, value *string, readFrom io.Reader, writeTo io.Writer) {
	if *value != "" {
		label = fmt.Sprintf("%v [%v]", label, *value)
	}
	got := GetInput(label, readFrom, writeTo)
	if got == "" {
		return
	}
	*value = got
}

func GetInput(label string, readFrom io.Reader, writeTo io.Writer) string {
	//Wait for an enter
	reader := bufio.NewReader(readFrom)
	if label != "" {
		fmt.Fprintf(writeTo, "%v: ", label)
	}
	text, _ := reader.ReadString('\n')

	text = strings.TrimSuffix(text, "\n")

	if text == "" {
		return ""
	}

	fmt.Fprintln(writeTo, "")

	return text
}

// Returns the index of the chosen option.
func GetChoice(choices []string, current int, readFrom io.Reader, writeTo io.Writer) int {
	if len(choices) == 0 {
		panic("No choices to choose from")
	}

	//Print choices
	for i, choice := range choices {
		fmt.Fprintf(writeTo, "%v. %v", i+1, choice)
		if i == current {
			fmt.Fprint(writeTo, " (current)\n")
		} else {
			fmt.Fprintln(writeTo)
		}
	}
	fmt.Fprintln(writeTo)

	fmt.Fprintf(writeTo, "Choose [1-%v]: ", len(choices))

	for {
		var err error = nil
		choiceStr := GetInput("", readFrom, writeTo)
		if choiceStr == "" {
			return current
		}

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Fprintln(writeTo, "Please enter a number")
		} else if choice < 1 {
			fmt.Fprintln(writeTo, "Number should be greater than 0")
		} else if choice > len(choices) {
			fmt.Fprintf(writeTo, "Number should not be larger than %v\n", len(choices))
		} else {
			return choice - 1
		}
	}
}
