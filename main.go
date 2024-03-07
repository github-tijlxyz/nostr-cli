package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

type Index struct {
	ActiveProfile string `yaml:"activeprofile"`
}

type Profile struct {
	Key struct {
		Public     string `yaml:"public"`
		Encryption string `yaml:"encryption"`
		Private    string `yaml:"private"`
	} `yaml:"key"`
	Metadata map[string][]map[string]interface{} `yaml:"metadata"`
}

var (
	index         Index
	openedProfile string
	s             Profile
)

func main() {
	app := &cli.App{
		Name:        "nostr-cli",
		Description: "some nostr cli tools",
		Version:     "0.3.0",
		Before: func(ctx *cli.Context) error {

			directory := ctx.String("directory")

			// Load index (for active profile)
			indexFile, err := os.ReadFile(filepath.Join(directory, "index.yaml"))
			if !errors.Is(err, os.ErrNotExist) {

				if err != nil {
					return err
				}
				err = yaml.Unmarshal(indexFile, &index)
				if err != nil {
					return err
				}
			} else {
				index = Index{
					ActiveProfile: "main",
				}
			}

			if index.ActiveProfile == "" {
				index.ActiveProfile = "main"
			}

			// Load config
			openedProfile = index.ActiveProfile
			configFile, err := os.ReadFile(filepath.Join(directory, fmt.Sprintf("%s.profile.yaml", index.ActiveProfile)))
			if !errors.Is(err, os.ErrNotExist) {
				if err != nil {
					return err
				}
				err = yaml.Unmarshal(configFile, &s)
				if err != nil {
					return err
				}
			} else {
				s = Profile{}
			}

			return nil
		},
		After: func(ctx *cli.Context) error {

			directory := ctx.String("directory")

			// Save index (for active profile)
			indexYaml, err := yaml.Marshal(&index)
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(directory, "index.yaml"), indexYaml, 0644)
			if err != nil {
				return err
			}

			// Save config
			if openedProfile == index.ActiveProfile {
				configYaml, err := yaml.Marshal(&s)
				if err != nil {
					return err
				}
				err = os.WriteFile(filepath.Join(directory, fmt.Sprintf("%s.profile.yaml", index.ActiveProfile)), configYaml, 0644)
				if err != nil {
					return err
				}
			}

			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "directory",
				Aliases: []string{"d"},
				Value:   getDefaultPath(),
				Usage:   "directory to save profile files and other data",
				EnvVars: []string{"NOSTR_CLI_DIRECTORY"},
			},
		},
		Commands: []*cli.Command{
			&cli.Command{
				Name:    "event",
				Aliases: []string{"ev"},
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:      "sign",
						Args:      true,
						ArgsUsage: "event json",
						Action: func(ctx *cli.Context) error {
							var event nostr.Event
							err := json.Unmarshal([]byte(ctx.Args().First()), &event)
							if err != nil {
								return err
							}
							signed, err := signEvent(event)
							if err != nil {
								return err
							}
							signedString, err := json.Marshal(signed)
							if err != nil {
								return err
							}

							fmt.Println(signedString)

							return nil
						},
					},
					&cli.Command{
						Name:      "publish",
						Args:      true,
						ArgsUsage: "event json",
						Action: func(ctx *cli.Context) error {
							var event nostr.Event
							err := json.Unmarshal([]byte(ctx.Args().First()), &event)
							if err != nil {
								return err
							}
							if event.Sig == "" {
								event, err = signEvent(event)
								if err != nil {
									return err
								}
							}

							err = publishEvent(event, []string{})

							return nil
						},
					},
					&cli.Command{
						Name:      "verify",
						Args:      true,
						ArgsUsage: "event json",
						Action: func(ctx *cli.Context) error {
							var event nostr.Event
							err := json.Unmarshal([]byte(ctx.Args().First()), &event)
							if err != nil {
								return err
							}
							verified, err := event.CheckSignature()
							if err != nil {
								return err
							}
							if verified {
								fmt.Println("valid")
							} else {
								fmt.Println("invalid")
							}
							return nil
						},
					},
				},
			},
			&cli.Command{
				Name: "key",
				Subcommands: []*cli.Command{
					&cli.Command{
						Name: "set",
					},
					&cli.Command{
						Name: "generate",
						Action: func(ctx *cli.Context) error {
							sk := nostr.GeneratePrivateKey()
							nsec, err := nip19.EncodePrivateKey(sk)
							pk, err := nostr.GetPublicKey(sk)
							npub, err := nip19.EncodePublicKey(pk)
							if err != nil {
								return err
							}

							fmt.Printf("\nnpub: %s\nnsec: %s\n\n", npub, nsec)

							set := confirm(fmt.Sprintf("do you want to set this keypair for profile '%s'", index.ActiveProfile), false)
							if set {
								err = setKey(nsec)
								if err != nil {
									return err
								}
							}

							return nil
						},
					},
				},
			},
			&cli.Command{
				Name:      "profile",
				Args:      true,
				ArgsUsage: "change to profile",
				Action: func(ctx *cli.Context) error {
					if ctx.Args().First() != "" {
						index.ActiveProfile = ctx.Args().First()
					}
					fmt.Println(index.ActiveProfile)
					return nil
				},
			},
			&cli.Command{
				Name: "profiles",
                Action: func(ctx *cli.Context) error {
                    directory := ctx.String("directory")
                    files, err := os.ReadDir(directory)
                    if err != nil {
                        return err
                    }
                    for _, file := range files {
                        if strings.HasSuffix(file.Name(), ".profile.yaml") {
                            fmt.Println(strings.TrimSuffix(file.Name(), ".profile.yaml"))
                        }
                    }

                    return nil
                },
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln("error:", err)
	}
}
