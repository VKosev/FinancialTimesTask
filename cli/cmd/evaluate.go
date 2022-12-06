package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vkosev/ft_cli/expression"
)

// evaluateCmd represents the evaluate command
var evaluateCmd = &cobra.Command{
	Use:   "evaluate",
	Short: "Evaluates an arithmetic expression.",
	Long:  `Evaluates an arithmetic expression based on text. Display an error if any occurs`,
	Run: func(cmd *cobra.Command, args []string) {
		client := expression.NewClient(nil)

		expr, err := cmd.Flags().GetString("expression")
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
			return
		}

		result, err := client.Evaluate(expr)
		if err != nil {
			if errors.Is(err, expression.ErrEvaluation) {
				fmt.Printf("\nERROR:\n")
				fmt.Println("  - Message: ", result.Message)
				fmt.Println("  - Type:    ", result.Type)

				return
			}
			fmt.Println("ERROR: ", err.Error())
			return
		}

		fmt.Printf("\nRESULT: %d", result.Result)
	},
}

func init() {
	rootCmd.AddCommand(evaluateCmd)

	evaluateCmd.Flags().StringP("expression", "e", "", "The expression to evaluate")
}
