package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/mdp/qrterminal/v3"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
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
    Use: "view",
    Short: "View the currently used key",
    Run: func(cmd *cobra.Command, args []string) {
        var (
            nsec string = ""
            npub string = ""
            privateKey string = ""
            publicKey string = ""
            err error
        )

        if viewKeyShowPrivate == true {
            privateKey, err = getKey()
            if err != nil {
                fmt.Println("error while getting key:", err)
                return
            }
            if privateKey == "" {
                fmt.Println("key not set")
            }
            publicKey, err = nostr.GetPublicKey(privateKey)
            if err != nil {
                fmt.Println("error while getting key:", err)
                return
            }
        } else {
            publicKey = viper.GetString("key.public")
            if publicKey == "" {
                fmt.Println("key not set")
            }
        }


        if publicKey != "" {
            npub, err = nip19.EncodePublicKey(publicKey)
            if err != nil {
                fmt.Println("error encoding npub:", err)
                return
            }
            fmt.Println("public key:", publicKey)
            fmt.Println("npub:", npub)
            if viewKeyViewQR == true {
                fmt.Print("\n")
                qrterminal.Generate(npub, qrterminal.L, os.Stdout)
                fmt.Print("\n")
            }
        }
        if privateKey != "" {
            nsec, err = nip19.EncodePrivateKey(privateKey)
            if err != nil {
                fmt.Println("error encoding nsec:", err)
                return
            }
            fmt.Println("private key:", privateKey)
            fmt.Println("nsec:", nsec)
        }

    },
}

var setKeyCmd = &cobra.Command{
    Use: "set",
    Short: "set new private key",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Print("enter new private key: ")
        inputBytes, err := term.ReadPassword(int(syscall.Stdin))
        input := strings.TrimSpace(string(inputBytes))
        if err != nil {
            fmt.Println("error while reading input:", err)
            return
        }
        fmt.Print("\n")
        setKey(input)
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

    publicKey, err := nostr.GetPublicKey(key)
    if err != nil {
        fmt.Println("error getting public key:", err)
        return
    }

    fmt.Print("Do you want to remove existing metadata? [y/N]: ")
    var input string = ""
    _, err = fmt.Scanln(&input)
    /*if err != nil {
        fmt.Println("error getting input:", err)
        return
    }*/
    input = strings.TrimSpace(strings.ToLower(input))
    switch input {
    case "y", "yes":
        clearMeta := &UserArgs{}
        viper.Set("metadata", clearMeta)
    default:
        //return
    }

    // encrypt the key?
    fmt.Print("\nDo you want to use encryption to store the private key? [Y/n]: ")
    input = ""
    _, err = fmt.Scanln(&input)
    /*if err != nil {
        fmt.Println("error getting input:", err)
        return
    }*/
    input = strings.TrimSpace(strings.ToLower(input))
    switch input {
    case "n", "no":
        viper.Set("key.public", publicKey)
        viper.Set("key.encryption", "none")
        viper.Set("key.secret", key)
    default:
        fmt.Print("\nEnter encryption passpharse: ")
        bytePassword, err := term.ReadPassword(int(syscall.Stdin))
        password := strings.TrimSpace(string(bytePassword))

        fmt.Print("\nConfirm encryption passpharse: ")
        byteConfirmPassword, err := term.ReadPassword(int(syscall.Stdin))
        confirmPassword := strings.TrimSpace(string(byteConfirmPassword))

        if err != nil {
            fmt.Println("error reading password:", err)
            return
        }

        fmt.Print("\n")

        if password != confirmPassword {
            fmt.Println("did not match")
            return
        }

        cipherKey, err := encrypt(key, password)
        if err != nil {
            fmt.Println("error encrypting: ", err)
            return
        }

        viper.Set("key.public", publicKey)
        viper.Set("key.encryption", "aes256+pbkdf2")
        viper.Set("key.secret", cipherKey)
    }

    err = viper.WriteConfig()
    if err != nil {
        fmt.Println("Failed to write config file:", err)
        return
    }

    fmt.Println("private key set")
}

func getKey () (string, error) {
    encryptionType := viper.GetString("key.encryption")

    if encryptionType == "none" {
        key := viper.GetString("key.secret")
        return key, nil
    } else if encryptionType == "aes256+pbkdf2" {

        cipherKey := viper.GetString("key.secret")

        fmt.Print("\nPlease enter your passpharse to decrypt the private key: ")
        passwordBytes, err := term.ReadPassword(int(syscall.Stdin))

        fmt.Print("\n\n")

        if err != nil {
            return "", err
        }

        password := strings.TrimSpace(string(passwordBytes))

        key, err := decrypt(cipherKey, password)
        if err != nil {
            return "", err
        }

        _, err = nip19.EncodePrivateKey(key)
        if err != nil {
            return "", errors.New("probably wrong password")
        }

        return key, nil

    } else {
        return "", errors.New("key not key or invalid encryption")
    }
}

