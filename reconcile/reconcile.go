package reconcile

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

var app *cli.App

func init() {
	if app != nil {
		return
	}
	app = &cli.App{
		Name:  "reconcile",
		Usage: "performs the reconciliation process by comparing transactions within the specified timeframe across system and bank statement data.",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Transactions processed\t:", 0)
			fmt.Println("Matched transactions\t:", 0)
			fmt.Println("Unmatched transactions\t:", 0)
			fmt.Println()
			fmt.Println("Missing bank statement:")
			for i, v := range []string{} {
				fmt.Printf("%d. %s\n", i+1, v)
			}
			fmt.Println()
			for _, file := range []string{} {
				fmt.Printf("%s - Missing system transaction:\n", file)
				// TODO placeholder result from using file parameter
				for i, v := range []string{} {
					fmt.Printf("%d. %s\n", i+1, v)
				}
				fmt.Println()
			}
			fmt.Println("Total disrepancy\t:", 0)
			return nil
		},
	}
}

func Reconcile(args []string) {
	if err := app.Run(args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
