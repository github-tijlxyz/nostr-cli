package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nbd-wtf/go-nostr"
)

func signEvent(event nostr.Event) (nostr.Event, error) {
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

func publishEvent(event nostr.Event, relays []string) error {
	if len(relays) < 1 {
		return errors.New("no relays too publish to")
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
