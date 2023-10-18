package main

import (
	"os"
    "fmt"
    "path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var customRelays string

var rootCmd = &cobra.Command{
    Use: "nostr-cli <command> [subcommand]",
    Short: "A command line interface for nostr",
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

var profileCmd = &cobra.Command{
    Use: "profile [command]",
    Run: profilePublishCmd.Run,
}

func init() {
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/nostr-cli/config.yaml)")
    //rootCmd.PersistentFlags().StringVar(&customRelays, "relays", "", "set relays (by default will use what is in config.yaml)")
    
    keyCmd.AddCommand(setKeyCmd)
    genKeyCmd.Flags().BoolVar(&genKeySet, "set", false, "directly set the generated key")
    genKeyCmd.Flags().BoolVar(&genKeyDontSet, "dont-set", false, "dont ask for setting the generated key")
    keyCmd.AddCommand(genKeyCmd)
    keyCmd.AddCommand(viewKeyCmd)
    viewKeyCmd.Flags().BoolVar(&viewKeyShowPrivate, "show-private", false, "show the private key")
    viewKeyCmd.Flags().BoolVar(&viewKeyViewQR, "qr", false, "print the npub as QR in terminal")
    rootCmd.AddCommand(keyCmd)

    eventCmd.AddCommand(signEventCmd)
    eventCmd.AddCommand(verifyEventCmd)
    eventCmd.AddCommand(publishEventCmd)
    rootCmd.AddCommand(eventCmd)

    rootCmd.AddCommand(profileCmd)

    relaysCmd.AddCommand(relaysSetCmd)
    relaysCmd.AddCommand(relaysViewCmd)
    relaysCmd.AddCommand(relaysAddCmd)
    relaysCmd.AddCommand(relaysRmCmd)
    rootCmd.AddCommand(relaysCmd)
}

func initConfig() error {
    if cfgFile == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            return err
        }
        cfgFile = filepath.Join(home, ".config", "nostr-cli", "config.yaml")
    }

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

    return nil
}

func main () {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
