package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
    log_n        int = 20
    cfgFolder    string
    profileName  string
    customRelays string
    relays       []string
)

var rootCmd = &cobra.Command {
    Use: "nostr-cli [command] [subcommand]",
    Short: "A command line interface for nostr",
    Version: "0.2.0",
    PersistentPreRunE: func (cmd *cobra.Command, args []string) error {
        return initConfig()
    },
}

var keyCmd = &cobra.Command{
    Use: "key [command]",
    Run: viewKeyCmd.Run,
}

var eventCmd = &cobra.Command{
    Use: "event <command>",
}

var relaysCmd = &cobra.Command{
    Use: "relays [command]",
    Run: relaysViewCmd.Run,
}

var metaCmd = &cobra.Command{
    Use: "metadata [command]",
    Run: profilePublishCmd.Run,
}

var feedCmd = &cobra.Command{
    Use: "feed [command]",
    Run: feedSubCmd.Run,
}

var configCmd = &cobra.Command{
    Use: "config",
    Run: func(cmd *cobra.Command, args []string) {
        location := viper.ConfigFileUsed()
        
        editor := os.Getenv("EDITOR")
        if editor == "" {
            editor = os.Getenv("VISUAL")
        }
        if editor == "" {
            editor = "vi"
        }
        ecmd := exec.Command(editor, location)
        ecmd.Stdin = os.Stdin
        ecmd.Stdout = os.Stdout
        ecmd.Stderr = os.Stderr

        if err := ecmd.Run(); err != nil {
            fmt.Println("error while trying to open editor:", err)
            return
        }
    },
}

func init() {
    rootCmd.PersistentFlags().StringVarP(&cfgFolder, "directory", "d", "", "directory to store config files. default is \"$HOME/.config/nostr-cli\". You can also set with the environment variable NOSTR_CLI_DIRECTORY")
    rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "p", "", "profile default is \"main\". you can also set the profile with the environment variable NOSTR_CLI_PROFILE")
    rootCmd.PersistentFlags().StringVarP(&customRelays, "relays", "r", "", "use relays (by default will use what is in the config)")

    setKeyCmd.Flags().IntVar(&log_n, "log-n", 22, "number of encryption rounds, set this lower for less powerful devices. set between 16 and 22")
    keyCmd.AddCommand(setKeyCmd)
    genKeyCmd.Flags().IntVar(&log_n, "log-n", 22, "number of encryption rounds, set this lower for less powerful devices. set between 16 and 22")
    genKeyCmd.Flags().BoolVar(&genKeySet, "set", false, "directly set the generated key")
    genKeyCmd.Flags().BoolVar(&genKeyDontSet, "dont-set", false, "dont ask for setting the generated key")
    keyCmd.AddCommand(genKeyCmd)
    keyCmd.AddCommand(viewKeyCmd)
    viewKeyCmd.Flags().BoolVar(&viewKeyShowPrivate, "private", false, "show the private key")
    viewKeyCmd.Flags().BoolVar(&viewKeyViewQR, "qr", false, "print the npub as QR in terminal")
    rootCmd.AddCommand(keyCmd)

    eventCmd.AddCommand(signEventCmd)
    eventCmd.AddCommand(verifyEventCmd)
    eventCmd.AddCommand(publishEventCmd)
    rootCmd.AddCommand(eventCmd)

    rootCmd.AddCommand(metaCmd)

    relaysCmd.AddCommand(relaysSetCmd)
    relaysCmd.AddCommand(relaysViewCmd)
    relaysCmd.AddCommand(relaysAddCmd)
    relaysCmd.AddCommand(relaysRmCmd)
    rootCmd.AddCommand(relaysCmd)

    rootCmd.AddCommand(configCmd)

    rootCmd.AddCommand(feedCmd)
}

func initConfig() error {

    if profileName == "" {
        profileName = "main"
    }

    if cfgFolder == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            return err
        }

        cfgFolder = filepath.Join(home, ".config", "nostr-cli")
    }

    cfgFile := filepath.Join(cfgFolder, fmt.Sprintf("%s.config.yaml", profileName))

    configDir := filepath.Dir(cfgFile)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		defaultConfig := []byte{} 
		if err := os.WriteFile(cfgFile, defaultConfig, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
    }

    viper.SetConfigFile(cfgFile)

    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("failed to read config file: %v\n", err)
    }

    if customRelays == "" {
        relays = viper.GetStringSlice("relays")
    } else {
        relays = strings.Split(customRelays, ",")
    }

    return nil
}

func main () {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

