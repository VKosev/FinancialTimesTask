package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vkosev/ft_cli/expression"
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate expression",
	Long:  `Check wether an expressions is valid`,
	Run: func(cmd *cobra.Command, args []string) {
		client := expression.NewClient(nil)

		expr, err := cmd.Flags().GetString("expression")
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return
		}

		fmt.Printf("\nStarting to validate expression %s\n", expr)
		result, err := client.Validate(expr)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return
		}

		if !result.Valid {
			fmt.Printf("\nValidation Result\n")
			fmt.Println("  - Valid:  ", result.Valid)
			fmt.Println("  - Reason: ", result.Reason)
		} else {
			fmt.Printf("\nValidation Result\n")
			fmt.Println("  - Valid: ", result.Valid)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	validateCmd.Flags().StringP("expression", "e", "", "the expression to evaluate")
}
