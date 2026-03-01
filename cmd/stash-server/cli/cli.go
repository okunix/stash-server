package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"gitlab.com/stash-password-manager/stash-server/cmd/stash-server/app"
	"gitlab.com/stash-password-manager/stash-server/version"
)

var rootCmd = &cobra.Command{
	Use:   "stash-server",
	Short: "Headless password manager",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		runCmd.Run(cmd, args)
	},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath, _ := cmd.Flags().GetString("config")
		app.Run(configFilePath)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show stash-server version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s\n", version.Version())
	},
}

func init() {
	rootCmd.PersistentFlags().
		String("config", "/etc/stash/server.yaml", "server config file location")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
