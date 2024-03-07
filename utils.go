package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func padString(input string, length int) string {
	if len(input) >= length {
		return input
	}
	return input + " " + strings.Repeat(" ", length-len(input)-1)
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getDefaultPath() string {

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".config", "nostr-cli")

}

func confirm(prompt string, defaultValue bool) bool {
	reader := bufio.NewReader(os.Stdin)

	var s string

	if defaultValue {
		fmt.Printf("%s [Y/n]: ", prompt)
	} else {
		fmt.Printf("%s [y/N]: ", prompt)
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(strings.ToLower(input[:len(input)-1]))

	switch s {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		// If no valid input is provided, use the default value.
		return defaultValue
	}
}
