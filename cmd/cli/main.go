package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/urfave/cli/v3"
)

func main() {
	setLogger()

	cmd := &cli.Command{
		Name:        "n8n-lint",
		Aliases:     nil,
		Usage:       "Simple n8n workflow JSON linter.",
		Version:     "",
		Description: "Simple lint tool for n8n workflow JSON files.",
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "check n8n workflow file(s) using a glob pattern",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					files, err := filepath.Glob(cmd.Args().First())
					if err != nil {
						return err
					}

					log.Println("found files:", len(files))
					log.Println(files)
					return nil
				},
			},
		},
	}

	log.Println("starting n8n-lint")

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func setLogger() {
	log.SetPrefix(chalk.Cyan("n8n-lint "))

	if val, ok := os.LookupEnv("ENV"); ok {
		if val == "DEV" {
			log.Println(chalk.Cyan("DEVELOPMENT MODE"))
			log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmsgprefix)
			return
		}
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmsgprefix)

}
