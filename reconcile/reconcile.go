package reconcile

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var app *cli.App
var param struct {
	SourceFile     string
	StatementFiles []string
	Start          cli.Timestamp
	End            cli.Timestamp
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
			&cli.BoolFlag{
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
			},
		},
		Action: func(ctx *cli.Context) error {
			fmt.Println("Transactions processed\t:", 0)
			fmt.Println("Matched transactions\t:", 0)
			fmt.Println("Unmatched transactions\t:", 0)
			for i, v := range []string{} {
				// 1. txID	filename.csv - uniqueID: -18.245,55
				fmt.Printf("%d. %s\t%s - %s: %.2f\n", i+1, v, v, v, 0.01)
			}
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
		Before: func(ctx *cli.Context) error {
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
