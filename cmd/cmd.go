package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cdb/actions-parser/expressions"
	"github.com/spf13/cobra"
)

var (
	runContext string
	rootCmd    = &cobra.Command{
		Use:   "actions-parser",
		Short: "An actions expression parser",
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runContext, "context", "c", "", "json formatted context")
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

		var jsonContext expressions.Context
		if runContext != "" {
			json.Unmarshal([]byte(runContext), &jsonContext)
			fmt.Println("Including context", jsonContext)
		}

		ast := expressions.Parse(args[0])
		out, err := expressions.Evaluate(ast, jsonContext)
		if err != nil {
			fmt.Println("An error occurred:", err)
		}
		fmt.Println("Output:  ", out)
	},
}
