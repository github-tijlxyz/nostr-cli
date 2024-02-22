package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	
	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
)


var feedSubCmd = &cobra.Command{
    Short: "sub <'json'>",
    Use: "subscribe live to filter",
    Run: func(cmd *cobra.Command, args []string) {
        var filter nostr.Filter
        err := json.Unmarshal([]byte(args[0]), &filter)
        if err != nil {
            fmt.Println("error decoding arg:", err)
            return
        }

        gotEvents := map[string]bool{}
        var mutex sync.Mutex

        ctx := context.Background()

        for _, url := range relays {
            go func(url string) {
                relay, err := nostr.RelayConnect(ctx, url)
                if err != nil {
                    fmt.Printf("\nerror while connecting to %s: %s", url, err)
                    return
                }

                sub, err := relay.Subscribe(ctx, []nostr.Filter{filter})
                if err != nil {
                    fmt.Printf("\nerror while subscribing to %s: %s", url, err)
                }

                fmt.Println("\nconnected to", url)
                for ev := range sub.Events {
                    mutex.Lock()
                    if gotEvents[ev.ID] {
                        mutex.Unlock()
                        continue
                    }
                    gotEvents[ev.ID] = true
                    mutex.Unlock()
                    fmt.Printf("\ngot event %s from %s", ev.ID, url)
                    fmt.Printf("\nkind: %v", ev.Kind)
                    fmt.Printf("\ntags: %v", ev.Tags)
                    fmt.Printf("\ncontent: %v", ev.Content)
                    fmt.Print("\n---")
                }
            }(url)
        }
        select {}
    },
}

