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
	Short: "mop is a tool to support manual operations using GitHub Issues",
	Long: `mop is a tool to support manual operations using GitHub Issues for operation definition management.

To use this tool, first prepare a GitHub repository and create an Issue. 
Next, write repository information to a configuration file (~/.mop.json). 
To start an operation, execute "mop operation" command.
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
