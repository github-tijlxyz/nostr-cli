package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
    "errors"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
)

func signEvent (event nostr.Event) (nostr.Event, error) {
    sk, err := getKey()
    if err != nil {
        return event, err
    }

    if event.CreatedAt == 0 {
        event.CreatedAt = nostr.Timestamp(time.Now().Unix())
    }

    err = event.Sign(sk)
    if err != nil {
        return event, err
    }

    return event, nil
}

func publishEvent (event nostr.Event, relays []string) error {
    if len(relays) < 1 {
        return errors.New("no relays set")
    }
    fmt.Printf("\nPublishing event %s to relays!\n", event.ID)
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

        err = relay.Publish(ctx, event)
        if err != nil {
            m := fmt.Sprintf("error while publishing to %v: %v", url, err)
            fmt.Printf("\r%s\n", padString(m, len(messageToReplace)))
            continue
        }

        m := fmt.Sprintf("published to %s", url)
        fmt.Printf("\r%s\n", padString(m, len(messageToReplace)))
    }

    return nil
}

var verifyEventCmd = &cobra.Command{
    Use: "verify <'event json'>",
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
    Use: "sign <'event json'>",
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
            fmt.Println("error signing event:", err)
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
    Use: "publish <'event json'>",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        arg := args[0]
        
        var event nostr.Event
        err := json.Unmarshal([]byte(arg), &event)
        if err != nil {
            fmt.Println("error reading event:", err)
            return
        }

        if event.Sig == "" {
            event, err = signEvent(event)
            if err != nil {
                fmt.Println("error signing event:", err)
                return
            }
        }

        // publish the event!
        err = publishEvent(event, relays)
        if err != nil {
            fmt.Println("error while publishing event:", err)
        }
    },
}

