package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mdp/qrterminal/v3"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var genKeySet bool
var genKeyDontSet bool
var viewKeyViewQR bool 
var viewKeyShowPrivate bool

var genKeyCmd = &cobra.Command{
    Use: "generate",
    Short: "Generate a new nostr keypair",   
    Run: func(cmd *cobra.Command, args []string) {
        sk := nostr.GeneratePrivateKey()
        pk, _ := nostr.GetPublicKey(sk)
        nsec, _ := nip19.EncodePrivateKey(sk)
        npub, _ := nip19.EncodePublicKey(pk)

        fmt.Print("\n")
        fmt.Println("Generated Key:")
        fmt.Println("private key:", sk)
        fmt.Println("public key:", pk)
        fmt.Println("nsec:", nsec)
        fmt.Println("npub:", npub)

        if genKeyDontSet {
            return
        }

        if !genKeySet {
            // set it as input?
            fmt.Print("\nDo you want to set this nostr keypair as private key for nostr-cli? [y/N]: ")
            var input string
            _, err := fmt.Scanln(&input)
            if err != nil {
                return
            }
            input = strings.TrimSpace(strings.ToLower(input))
            switch input {
            case "y", "yes":
                fmt.Print("\n")
                setKey(sk)
            case "n", "no":
                return
            default:
                return
            }
        } else if genKeySet {
            fmt.Print("\n")
            setKey(sk)
        }


    },
}

var viewKeyCmd = &cobra.Command{
    Use: "view [encoded key]",
    Args: cobra.RangeArgs(0, 1),
    Short: "View the currently used keypair, or convert a nip19 encoded key to hex",
    Run: func(cmd *cobra.Command, args []string) {
        if len(args) > 0 {
            viewKey(args[0], true)
        } else {
            key := viper.GetString("key")
            if key == "" {
                fmt.Println("key not set")
                return
            }
            viewKey(key, false)
        }
    },
}

func viewKey (key string, nip19convert bool) {
    var (
        nsec string
        npub string
        sk string
        pk string

        err error
    )
    if strings.HasPrefix(key, "nsec1") {
        var vsk any
        _, vsk, err = nip19.Decode(key)
        if err != nil {
            fmt.Println("error decoding nsec:", err)
            return
        }
        sk = vsk.(string)
        pk, err = nostr.GetPublicKey(sk)
        npub, err = nip19.EncodePublicKey(pk)
        nsec, err = nip19.EncodePrivateKey(sk)
    } else if strings.HasPrefix(key, "npub1") {
        var vpk any
        _, vpk, err = nip19.Decode(key)
        if err != nil {
            fmt.Println("error decoding npub:", err)
            return
        }
        pk = vpk.(string)
        npub, err = nip19.EncodePublicKey(pk)
    } else if nip19convert == false {
        _, terr := nostr.GetPublicKey(key)
        if terr != nil { // is it a pk?
            pk = key
            npub, err = nip19.EncodePublicKey(pk)
        } else { // is it a sk?
            sk = key
            nsec, err = nip19.EncodePrivateKey(sk)
            pk, err = nostr.GetPublicKey(sk)
            npub, err = nip19.EncodePublicKey(pk)
        }
    } else {
        fmt.Println("please use a nip19 encoded key")
        return
    }
    if err != nil {
        fmt.Println("something went wrong", err)
        return
    }
    if viewKeyViewQR == true {
        fmt.Print("\n")
        qrterminal.Generate(npub, qrterminal.L, os.Stdout)
        fmt.Print("\n")
    } else {
        if sk != "" {
            fmt.Println("pk:", pk)
        }
        if nsec != "" {
            fmt.Println("npub:", npub)
        }
        if nip19convert == true || viewKeyShowPrivate == true {
            if pk != "" {
                fmt.Println("sk:", sk)
            }
            if npub != "" {
                fmt.Println("nsec:", nsec)
            }
        }
    }
}

var setKeyCmd = &cobra.Command{
    Use: "set <private key>",
    Short: "Set private key",
    Args: cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        setKey(args[0])
    },
}

func setKey (arg string) {
    key := arg

    if strings.HasPrefix(key, "nsec1") {
        // decode nsec
        _, decoded, err := nip19.Decode(key)
        if err != nil {
            fmt.Println("Invalid Key:", err)
            return
        }
        key = decoded.(string)
    } else {
        // check if we can encode to nsec
        _, err := nip19.EncodePrivateKey(key)
        if err != nil {
            fmt.Println("Invalid Key:", err)
            return
        }
    }

    viper.Set("key", key)

    err := viper.WriteConfig()
    if err != nil {
        fmt.Println("Failed to write config file:", err)
        return
    }

    fmt.Println("private key set")
}

//func accessKey () string {
// first decrypt private key or something
//}

