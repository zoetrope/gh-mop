package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/config"
)

var configPath string
var mopConfig *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mop",
	Short: "mop is a tool for managing manual operations",
	Long: `mop is a tool for managing manual operations using GitHub issues as a backend.

At first, prepare a GitHub repository and create an issue.
Then, prepare a configuration file and run the "mop-start" command.
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if configPath == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configPath = home + "/.mop.json"
		}
		mopConfig, err = config.LoadConfig(configPath)
		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.mop.json)")
}
