package main

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/nbd-wtf/go-nostr/nip49"
	"golang.org/x/term"
)

func setKey(nsec string) error {
	var err error
	if !strings.HasPrefix(nsec, "nsec1") {
		nsec, err = nip19.EncodePrivateKey(nsec)
		if err != nil {
			log.Println("a")
			return err
		}
	}

	prefix, privateKeyAny, err := nip19.Decode(nsec)
	if err != nil {

		log.Println("b")
		return err
	}
	if prefix != "nsec" {

		log.Println("c")
		return errors.New("not a nsec")
	}
	privateKey := privateKeyAny.(string)

	publicKey, err := nostr.GetPublicKey(privateKey)
	if err != nil {
		log.Println("d")
		return err
	}
	npub, err := nip19.EncodePublicKey(publicKey)
	if err != nil {
		log.Println("e")
		return err
	}

	if len(s.Metadata) != 0 {
		c := confirm(fmt.Sprintf("do you want to remove existing metadata for profile '%s'?", index.ActiveProfile), false)
		if c {
			s.Metadata = make(map[string][]map[string]interface{})
		}
	}

	c := confirm("do you want to encrypt the private key?", true)
	if !c {
		s.Key.Encryption = "none"
		s.Key.Private = nsec
		s.Key.Public = npub
	} else {
		fmt.Printf("\nNew encryption passpharse: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		password := strings.TrimSpace(string(bytePassword))

		fmt.Printf("\nConfirm encryption passpharse: ")
		byteConfirmPassword, err := term.ReadPassword(int(syscall.Stdin))
		confirmPassword := strings.TrimSpace(string(byteConfirmPassword))

		if err != nil {
			return err
		}

		if password != confirmPassword {
			return errors.New("passwords did not match")
		}

		fmt.Println()

		cipherKey, err := nip49.Encrypt(privateKey, password, uint8(21), nip49.ClientDoesNotTrackThisData)
		if err != nil {
			return err
		}

		s.Key.Encryption = "nip49"
		s.Key.Public = npub
		s.Key.Private = cipherKey
	}

	return nil
}

func getKey() (string, error) {
	if s.Key.Public == "" || s.Key.Private == "" {
		return "", errors.New("key not set")
	}
	if s.Key.Encryption == "nip49" {
		fmt.Printf("\nPlease enter the passpharse to decrypt the private key:")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		if err != nil {
			return "", err
		}
		password := strings.TrimSpace(string(passwordBytes))
		key, err := nip49.Decrypt(s.Key.Private, password)
		if err != nil {
			return "", err
		}
		_, err = nip19.EncodePrivateKey(key)
		if err != nil {
			return "", errors.Join(errors.New("HINT: probably wrong password"), err)
		}
		return key, nil
	} else if s.Key.Encryption == "none" {
		if strings.HasPrefix(s.Key.Private, "nsec") {
			prefix, privateKey, err := nip19.Decode(s.Key.Private)
			if err != nil {
				return "", err
			}
			if prefix != "nsec" {
				return "", errors.New("not a nsec")
			}
			return privateKey.(string), nil
		}
		return s.Key.Private, nil
	}
	return "", errors.New("invalid encryption method")
}
