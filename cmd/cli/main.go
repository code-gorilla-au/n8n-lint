package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/urfave/cli/v3"
)

var Version = "dev-v1.0.0"

func main() {
	setLogger()

	cwd, cErr := os.Getwd()
	if cErr != nil {
		log.Fatal(cErr)
	}

	defaultConfigPath := filepath.Clean(filepath.Join(cwd, ".n8n-lint.yaml"))
	flagConfigPath := defaultConfigPath

	cmd := &cli.Command{
		Name:        "n8n-lint",
		Aliases:     nil,
		Usage:       "Simple n8n workflow JSON linter.",
		Version:     Version,
		Description: "Simple lint tool for n8n workflow JSON files.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "override the default config file path",
				Value:       defaultConfigPath,
				Destination: &flagConfigPath,
				Aliases:     []string{"c"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "check n8n workflow file(s) using a glob pattern",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					files, err := filepath.Glob(cmd.Args().First())
					if err != nil {
						return err
					}

					log.Println("config", defaultConfigPath)

					log.Println("found files:", len(files))
					log.Println(files)
					return nil
				},
			},
		},
	}

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
