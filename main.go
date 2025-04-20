package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var (
	HatoPatterns   []string = []string{"ぽ", "っ"}
	NormalPatterns []string = []string{"お", "あ"}
)

func getRandomFromArray(arr []string, size int) string {
	result := ""
	for range size {
		result += arr[rand.Intn(len(arr))]
	}
	return result
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "hato",
				Usage: "Set execution user as hato",
			},
			&cli.BoolFlag{
				Name:  "force",
				Usage: "Execute hakkyo forcefully",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "debug mode",
			},
		},
		Action: func(ctx *cli.Context) error {
			if ctx.Bool("debug") {
				slog.SetLogLoggerLevel(slog.LevelDebug)
			}
			isHato := ctx.Bool("hato")
			isForce := ctx.Bool("force")
			if isHato && !isForce {
				fmt.Println("鳩を発狂させることはできません。")
				os.Exit(1)
			}

			slog.Debug("Passed flags", slog.Bool("hato", isHato), slog.Bool("force", isForce))
			w, _, err := term.GetSize(int(os.Stdin.Fd()))
			if err != nil {
				slog.Debug("Failed to get terminal size. Fallback to default value (20)")
				w = 20
			}
			w--
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
		L:
			for {
				select {
				case <-sig:
					break L
				default:
					if isHato {
						fmt.Println(getRandomFromArray(HatoPatterns, w))
					} else {
						fmt.Println(getRandomFromArray(NormalPatterns, w))
					}
				}
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("Failed to execute hakkyo", slog.Any("error", err))
		os.Exit(1)
	}
}
