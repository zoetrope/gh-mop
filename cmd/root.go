package cmd

import (
	"fmt"
	"os"

	"github.com/cli/go-gh"
	"github.com/spf13/cobra"
	"github.com/zoetrope/gh-mop/config"
)

var configPath string
var mopConfig *config.Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mop",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hi world, this is the gh-mop extension!")
		client, err := gh.RESTClient(nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		response := struct{ Login string }{}
		err = client.Get("user", &response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("running as %s\n", response.Login)
	},
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.mop.json)")
}
