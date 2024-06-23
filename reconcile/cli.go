package reconcile

import (
	"fmt"

	"github.com/dewzzjr/amartha/reconcile/action"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var app *cli.App
var param action.Param
var debugFlag = &cli.BoolFlag{
	Name:    "verbose",
	Usage:   "Debug mode print all log",
	Aliases: []string{"v"},
	Action: func(ctx *cli.Context, debug bool) error {
		if debug {
			log.Logger = log.Level(zerolog.DebugLevel)
		} else {
			log.Logger = log.Level(zerolog.ErrorLevel)
		}
		return nil
	},
}

func init() {
	if app != nil {
		return
	}
	app = &cli.App{
		Name:  "reconcile",
		Usage: "performs the reconciliation process by comparing transactions within the specified timeframe across system and bank statement data.",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "source",
				Usage:       "System transaction CSV file path",
				Required:    true,
				Aliases:     []string{"file", "f"},
				Destination: &param.SourceFile,
			},
			&cli.MultiStringFlag{
				Target: &cli.StringSliceFlag{
					Name:      "statement",
					Usage:     "Bank statement CSV file path (can handle multiple files from different banks)",
					Required:  true,
					Aliases:   []string{"bank", "b"},
					TakesFile: true,
				},
				Destination: &param.StatementFiles,
			},
			&cli.TimestampFlag{
				Name:        "start",
				Usage:       "Start date for reconciliation timeframe. Format: yyyy-MM-dd",
				Layout:      "2006-01-02",
				Aliases:     []string{"s"},
				Destination: &param.Start,
			},
			&cli.TimestampFlag{
				Name:        "end",
				Usage:       "End date for reconciliation timeframe. Format: yyyy-MM-dd",
				Layout:      "2006-01-02",
				Aliases:     []string{"e"},
				Destination: &param.End,
			},
			debugFlag,
		},
		Action: func(ctx *cli.Context) error {
			result, err := action.Run(param)
			if err != nil {
				return err
			}
			fmt.Println("Transactions processed\t:", result.Total)
			fmt.Println("Matched transactions\t:", result.Matched)
			fmt.Println("Unmatched transactions\t:", result.Unmatched)
			fmt.Println()
			if len(result.TransactionMissing) > 0 {
				fmt.Printf("%s - Missing bank statement:\n", param.SourceFile)
				for i, v := range result.TransactionMissing {
					fmt.Printf("%d. %v\n", i+1, v)
				}
				fmt.Println()
			}
			for _, file := range param.StatementFiles {
				missing := result.BankMissing[file]
				if len(missing) == 0 {
					continue
				}
				fmt.Printf("%s - Missing system transaction:\n", file)
				for i, v := range missing {
					fmt.Printf("%d. %v\n", i+1, v)
				}
				fmt.Println()
			}
			// TODO calculate disrepancy
			fmt.Println("Total disrepancy\t:", 0)
			return nil
		},
		Before: func(ctx *cli.Context) error {
			debugFlag.RunAction(ctx)
			log.Debug().Any("param", param).Msg("input parameter")
			return nil
		},
		Suggest:              true,
		EnableBashCompletion: true,
	}
}

func Reconcile(args []string) {
	if err := app.Run(args); err != nil {
		log.Fatal().Msgf("%v", err)
	}
}
