package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var relaysSetCmd = &cobra.Command{
    Use: "set <ws://relay1,ws://relay2>",
    Short: "set the relays used by default",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        arg := args[0]
        array := strings.Split(arg, ",")
        viper.Set("relays", array)
        err := viper.WriteConfig()
        if err != nil {
            fmt.Println("Error writing config:", err)
            return
        }
    },
}

var relaysAddCmd = &cobra.Command{
    Use: "add <ws://relay>",
    Args: cobra.ExactArgs(1),
    Short: "Add a relay to use",
    Run: func(cmd *cobra.Command, args []string) {
        relays := viper.GetStringSlice("relays")
        relays = append(relays, args[0])
        viper.Set("relays", relays)
        err := viper.WriteConfig()
        if err != nil {
            fmt.Println("Error writing config:", err)
            return
        }
    },
}

var relaysRmCmd = &cobra.Command{
    Use: "rm <ws://relay>",
    Args: cobra.ExactArgs(1),
    Short: "Remove a relay to use",
    Run: func(cmd *cobra.Command, args []string) {
        relays := viper.GetStringSlice("relays")
        for i, str := range relays {
            if str == args[0] {
                copy(relays[i:], relays[i+1:])
                relays = relays[:len(relays)-1]
            }
        }
        viper.Set("relays", relays)
        err := viper.WriteConfig()
        if err != nil {
            fmt.Println("error writing config", err)
        }
    },
}

var relaysViewCmd = &cobra.Command{
    Use: "view",
    Short: "view currently used relays",
    Run: func(cmd *cobra.Command, args []string) {
        relays := viper.GetStringSlice("relays")
        if len(relays) < 1 {
            fmt.Println("no relays set")
            return
        }
        fmt.Println("currently used relays:")
        for i, v := range relays {
            fmt.Printf("\n%v: %v", i+1, v)
        }
    },
}


