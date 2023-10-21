package main

import (
	"os"
    "fmt"
    "path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
    "os/exec"
)

var cfgFile string
//var cfgProfile string
var customRelays string

var rootCmd = &cobra.Command{
    Use: "nostr-cli [command] [subcommand]",
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

var metaCmd = &cobra.Command{
    Use: "metadata [command]",
    Run: profilePublishCmd.Run,
}

var feedCmd = &cobra.Command{
    Use: "feed [command]",
    Run: feedSubCmd.Run,
}

/*var configprofileCmd = &cobra.Command{
    Use: "profile <config profile to activate>",
    Run: func(cmd *cobra.Command, args []string) {
        arg := args[0]
        if arg == "" {
            fmt.Println("missing args")
            return
        }
        viper.Set("active", arg)
        err := viper.WriteConfig()
        if err != nil {
            fmt.Println("error writing config:", err)
            return
        }
    },
}*/

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
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/nostr-cli/config.yaml)")
    //rootCmd.PersistentFlags().StringVar(&customRelays, "relays", "", "set relays (by default will use what is in config.yaml)")
    
    keyCmd.AddCommand(setKeyCmd)
    genKeyCmd.Flags().BoolVar(&genKeySet, "set", false, "directly set the generated key")
    genKeyCmd.Flags().BoolVar(&genKeyDontSet, "dont-set", false, "dont ask for setting the generated key")
    keyCmd.AddCommand(genKeyCmd)
    keyCmd.AddCommand(viewKeyCmd)
    viewKeyCmd.Flags().BoolVar(&viewKeyShowPrivate, "view-private", false, "show the private key")
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

    //rootCmd.AddCommand(configprofileCmd)

    rootCmd.AddCommand(feedCmd)
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

    /*cfgProfile = viper.GetString("active")
    if cfgProfile == "" {
        cfgProfile = "default"
    }*/

    return nil
}

func main () {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
