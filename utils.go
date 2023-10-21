package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func getPublicKey () string {
    publicKey := viper.GetString("key.public")
    if publicKey == "" {
        fmt.Println("key not set")
        return ""
    }
    return publicKey
}

func padString (input string, length int) string {
    if len(input) >= length {
        return input
    }
    return input + " " + strings.Repeat(" ", length-len(input)-1)
}

func getUserInput (prompt string) string {
    fmt.Print(prompt)
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}

func truncateString (s string) string {
    //l := 40
    //if len(s) > l {
    //    return s[:l-3] + "..."
    //}
    return s
}
