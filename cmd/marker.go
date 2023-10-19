/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// markerCmd represents the marker command
var markerCmd = &cobra.Command{
	Use:   "marker <new marker>",
	Short: "Set the marker for a todo comment",
	Long: `"Set the marker for a todo comment in the config file.
The default marker for todo comments is 'TODO:'.
Pass in a string that will be the marker.
If no string is passed in, will print the current marker.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("The current marker for todo comments is:", viper.GetString("marker"))
			return
		}

		marker := args[0]
		viper.Set("marker", marker)

		// Write the updated configuration to the config file
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error writing to config file: %v\n", err)
			return
		}

		fmt.Printf("Marker set to: %s\n", marker)

	},
}

func init() {
	rootCmd.AddCommand(markerCmd)
}
