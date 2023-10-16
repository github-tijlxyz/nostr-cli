package main

import (
	"encoding/json"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

        sk := viper.GetString("key")
        if sk == "" {
            fmt.Println("key not set")
            return
        }
        
        err = event.Sign(sk)
        if err != nil {
            fmt.Println("error signing event:", err)
            return
        }

        jsonBytes, err := json.Marshal(event)
        if err != nil {
            fmt.Println("error signing event:", err)
            return
        }

        fmt.Println("signed event json:")
        fmt.Println(string(jsonBytes))
    },
}

var publishEventCmd = &cobra.Command{
    Use: "publish <event json>",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
       // sign and publish event  
    },
}



