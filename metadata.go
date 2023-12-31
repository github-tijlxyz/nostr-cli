package main

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
    "github.com/inancgumus/screen"
)

type UserArgs struct {
    Name string `json:"name" yaml:"name"`
    Displayname string `json:"display_name" yaml:"displayname"`
    NIP05 string `json:"nip05" yaml:"nip05"`
    Picture string `json:"picture" yaml:"picture"`
    Banner string `json:"banner" yaml:"banner"`
    About string `json:"about" yaml:"about"`
    Website string `json:"website" yaml:"website"`
    LUD16 string `json:"lud16" yaml:"lud16"` // ln address
    LUD06 string `json:"lud06" yaml:"lud06"` // lnurl1
}

var profilePublishCmd = &cobra.Command{
    Use:   "set",
    Short: "set new profile data",
    Run: func(cmd *cobra.Command, args []string) {
        user := &UserArgs{}

        if err := viper.UnmarshalKey("metadata", user); err != nil {
            fmt.Println("error decoding userdata from config:", err)
            return
        }

        for {
            screen.Clear()
            displayMenu(user)
            option := getUserInput("Enter an option: ")

            switch option {
            case "q":
                return
            case "s", "w", "wq":
                // Save to config
                viper.Set("metadata", user)

                if err := viper.WriteConfig(); err != nil {
                    fmt.Println("error while writing config:", err)
                    return
                }

                if option == "w" || option == "wq" {
                    // Publish to relays
                    jsonData, err := json.Marshal(user)
                    if err != nil {
                        fmt.Println("error while encoding JSON:", err)
                        return
                    }

                    metaEvent := nostr.Event{
                        Content: string(jsonData),
                        Tags:    nostr.Tags{},
                        Kind:    0,
                    }

                    if metaEvent, err = signEvent(metaEvent); err != nil {
                        fmt.Println("error while signing event:", err)
                        return
                    }

                    err = publishEvent(metaEvent, relays)
                    if err != nil {
                        fmt.Println("error while publishing event:", err)
                        return
                    }

                    if option == "wq" {
                        return
                    }
                }
            case "l":
                if len(relays) < 1 {
                    fmt.Println("no relays set")
                    return
                }

                key := getPublicKey()
                filter := nostr.Filter{
                    Kinds:   []int{0},
                    Authors: []string{key},
                    Limit:   1,
                }

                e, err := getEventFromRelays(filter, relays)
                if err != nil {
                    fmt.Println("error while getting event from relays:", err)
                    return
                }

                if err = json.Unmarshal([]byte(e.Content), user); err != nil {
                    fmt.Println("error while decoding event content:", err)
                    return
                }
            default:
                updateField(option, user)
            }
        }
    },
}

func displayMenu(user *UserArgs) {
    fmt.Print("\nMenu:\n\n")
    t := reflect.TypeOf(*user)
    v := reflect.ValueOf(*user)

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := truncateString(v.Field(i).String())
        fmt.Printf("%d. %s: %s\n", i+1, field.Name, value)
    }

    fmt.Print("\n")
    fmt.Println("l. Load from relays (only do this if you published metadata from another client)")
    fmt.Println("s. Save to config")
    fmt.Println("w. Publish to relays and safe to config")
    fmt.Println("q. Quit (without saving)")    
    fmt.Println("wq. Publish to relays, safe to config and quit")
    fmt.Print("\n")

}

func updateField(option string, user *UserArgs) {
    index := int(option[0] - '1') // Convert the option to an index
    t := reflect.TypeOf(user).Elem()
    v := reflect.ValueOf(user).Elem()

    if index < 0 || index >= t.NumField() {
        fmt.Println("Invalid option. Please try again.")
        return
    }

    field := t.Field(index)
    fmt.Printf("Enter new value for %s: ", field.Name)
    newValue := getUserInput("")
    v.Field(index).SetString(newValue)
    fmt.Printf("Updated %s.\n", field.Name)
}

