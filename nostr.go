package main

import (
	"context"
    "time"
    "errors"

	"github.com/nbd-wtf/go-nostr"
)

func getEventFromRelays (filter nostr.Filter, relays []string) (nostr.Event, error) {

    ctx := context.Background()
    var resultEvents []nostr.Event 

    for _, url := range relays {
        relay, err := nostr.RelayConnect(ctx, url)
        if err != nil {
            continue
        }

        ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
        defer cancel()

        sub, err := relay.Subscribe(ctx, nostr.Filters{filter})
        if err != nil {
            continue
        }

        for ev := range sub.Events {
            resultEvents = append(resultEvents, *ev)
        }
    }

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
