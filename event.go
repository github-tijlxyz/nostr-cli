package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func signEvent (event nostr.Event) (nostr.Event, error) {
    sk := viper.GetString("key")
    if sk == "" {
        fmt.Println("key not set")
        return event, errors.New("key not set")
    }

    if event.CreatedAt == 0 {
        event.CreatedAt = nostr.Timestamp(time.Now().Unix())
    }

    err := event.Sign(sk)
    if err != nil {
        fmt.Println("error signing event:", err)
        return event, err
    }

    return event, nil
}

func publishEvent (event nostr.Event, relays []string) {
    ctx := context.Background()
    fmt.Print("\n")
    for _, url := range relays {
        messageToReplace := fmt.Sprintf("publishing to %s...", url)
        fmt.Print(messageToReplace)

        relay, err := nostr.RelayConnect(ctx, url)
        if err != nil {
            m := fmt.Sprintf("error while connecting to %v: %v", url, err)
            fmt.Printf("\r%s\n", padString(m, len(messageToReplace)))
            continue
        }

        _, err = relay.Publish(ctx, event)
        if err != nil {
            m := fmt.Sprintf("error while publishing to %v: %v", url, err)
            fmt.Printf("\r%s\n", padString(m, len(messageToReplace)))
            continue
        }

        m := fmt.Sprintf("published to %s", url)
        fmt.Printf("\r%s\n", padString(m, len(messageToReplace)))
    }
}

var verifyEventCmd = &cobra.Command{
    Use: "verify <event json>",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        arg := args[0]
        
        var event nostr.Event
        err := json.Unmarshal([]byte(arg), &event)
        if err != nil {
            fmt.Println("error reading event:", err)
            return
        }
        
        verified, err := event.CheckSignature()
        if err != nil {
            fmt.Println("error checking signature:", err)
            return
        }

        if verified {
            fmt.Println("valid signature")
        } else {
            fmt.Println("invalid signature")
        }
    },
}

var signEventCmd = &cobra.Command{
    Use: "sign <event json>",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        arg := args[0]
        
        var event nostr.Event
        err := json.Unmarshal([]byte(arg), &event)
        if err != nil {
            fmt.Println("error reading event:", err)
            return
        }

        event, err = signEvent(event)
        if err != nil {
            return
        }

        jsonBytes, err := json.Marshal(event)
        if err != nil {
            fmt.Println("error encoding signed event:", err)
            return
        }

        fmt.Println(string(jsonBytes))
    },
}

var publishEventCmd = &cobra.Command{
    Use: "publish <event json>",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        arg := args[0]
        
        var event nostr.Event
        err := json.Unmarshal([]byte(arg), &event)
        if err != nil {
            fmt.Println("error reading event:", err)
            return
        }

        event, err = signEvent(event)
        if err != nil {
            return
        }

        relays := viper.GetStringSlice("relays")

        if len(relays) < 1 {
            fmt.Println("no relays set")
            return
        }

        // publish the event!
        publishEvent(event, relays)
    },
}



