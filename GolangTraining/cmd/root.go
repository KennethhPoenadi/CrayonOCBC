package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var config string

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "My cool app",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetConfigFile(config)
		err := viper.ReadInConfig()
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world!")
	},
}

// Execute runs the root command.
func Execute() error {
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "config.yaml", "Config file to use")
	return rootCmd.Execute()
}
