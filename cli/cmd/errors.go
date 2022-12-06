/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vkosev/ft_cli/expression"
)

// errorsCmd represents the errors command
var errorsCmd = &cobra.Command{
	Use:   "errors",
	Short: "Retrieves all occurred errors of expression evaluation or validation.",
	Long:  `Retrieves all occurred errors of expression evaluation or validation.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := expression.NewClient(nil)

		errors, err := client.Errors()
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return
		}

		if len(errors) == 0 {
			fmt.Printf("\nThere were no registered errors.\n")
			return
		}

		for _, e := range errors {
			fmt.Printf("\nERROR HISTORY\n")
			fmt.Println("  - Expression: ", e.Expression)
			fmt.Println("  - Frequency:  ", e.Frequency)
			fmt.Println("  - Type:       ", e.ErrType)
			fmt.Println("  - Endpoints:  ")

			for _, endpoint := range e.Endpoints {
				fmt.Println("     - Url:   ", endpoint.Url)
				fmt.Println("     - Count: ", endpoint.Count)
				fmt.Println("")
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(errorsCmd)
}
