package main

import (
    "context"
    "errors"
    "sync"
    "time"
    "fmt"

    "github.com/nbd-wtf/go-nostr"
)

func getEventFromRelays (filter nostr.Filter, relays []string) (nostr.Event, error) {

    fmt.Println("\nGetting event from relays... (this can take a bit depending on how many relays you configured)")

    ctx := context.Background()
    var resultEvents []nostr.Event

    var wg sync.WaitGroup

    concurrencyLimit := 16
    taskCh := make(chan string, concurrencyLimit)

    executeTask := func(url string) {
        defer wg.Done()

        relay, err := nostr.RelayConnect(ctx, url)
        if err != nil {
            return
        }

        ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
        defer cancel()

        sub, err := relay.Subscribe(ctx, nostr.Filters{filter})
        if err != nil {
            return
        }

        for ev := range sub.Events {
            resultEvents = append(resultEvents, *ev)
        }
    }

    for _, url := range relays {
        wg.Add(1)
        taskCh <- url
        go func() {
            for t := range taskCh {
                executeTask(t)
            }
        }()
    }

    close(taskCh)

    wg.Wait()

    if len(resultEvents) < 1 {
        return nostr.Event{}, errors.New("could not find any events on relays")
    }

    var latestEvent nostr.Event

    for _, event := range resultEvents {
        if event.CreatedAt > latestEvent.CreatedAt {
            latestEvent = event
        }
    }

    return latestEvent, nil
}
