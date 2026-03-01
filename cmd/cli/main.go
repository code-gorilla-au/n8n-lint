package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/code-gorilla-au/n8n-lint/internal/chalk"
	"github.com/code-gorilla-au/n8n-lint/internal/engine"
	"github.com/code-gorilla-au/n8n-lint/internal/logging"
	"github.com/code-gorilla-au/n8n-lint/internal/n8n"
	"github.com/code-gorilla-au/n8n-lint/internal/rules"
	"github.com/urfave/cli/v3"
)

var Version = "dev"

func main() {
	setLogger()

	cwd, cErr := os.Getwd()
	if cErr != nil {
		log.Fatal(cErr)
	}

	defaultConfigPath := filepath.Clean(filepath.Join(cwd, ".n8n-lint.yaml"))
	flagConfigPath := defaultConfigPath
	flagVerbose := false

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
			&cli.BoolFlag{
				Name:        "verbose",
				Usage:       "enable verbose logging",
				Value:       false,
				Destination: &flagVerbose,
				Aliases:     []string{"v"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "check n8n workflow file(s) using a glob pattern",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if flagVerbose {
						logging.SetVerbose()
						logging.Log("verbose mode enabled")
					}

					config, err := rules.LoadConfigFromFile(flagConfigPath)
					if err != nil {
						return err
					}

					workflows, err := n8n.LoadWorkflowsFromDir(cwd, config.Include, config.Ignore)
					if err != nil {
						return err
					}

					orchestrator := engine.NewOrchestrator(config)

					summary, err := orchestrator.Run(workflows)
					if err != nil {
						return err
					}

					summary.Print()

					if summary.TotalErrors() > 0 {
						return fmt.Errorf("failed with %d errors", summary.TotalErrors())
					}

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
