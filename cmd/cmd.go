package cmd

import (
	"fmt"
	"os"

	"github.com/cdb/actions-parser/expressions"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "actions-parser",
		Short: "An actions expression parser",
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the passed in expression",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("You must pass an expression argument")
			os.Exit(1)
		}
		fmt.Println("Running: ", args[0])
		ast := expressions.Parse(args[0])
		out, err := expressions.Evaluate(ast, nil)
		if err != nil {
			fmt.Println("An error occurred:", err)
		}
		fmt.Println("Output:  ", out)
	},
}
