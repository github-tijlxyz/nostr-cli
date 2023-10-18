package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/viper"
)

func getPublicKey () string {
    sk := viper.GetString("key")
    pk, err := nostr.GetPublicKey(sk)
    if err != nil {
        fmt.Println("cannot get key from config")
        return ""
    }
    return pk
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
    l := 40
    if len(s) > l {
        return s[:l-3] + "..."
    }
    return s
}
