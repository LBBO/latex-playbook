package main

import (
	"context"
	"log"
	"os"

	"github.com/lbbo/latex-playbook/go-scripts/extractPageNumbers"
	occurrenceSubtitles "github.com/lbbo/latex-playbook/go-scripts/occurrence-subtitles"

	createSceneTableAction "github.com/lbbo/latex-playbook/go-scripts/sceneTable"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "latex playbook scripts",
		Usage: "Scripts for managing latex playbook",
		Commands: []*cli.Command{
			{
				Name:  "create-scene-table",
				Usage: "Create scene table",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "src",
						Usage:     "Path to the latex source folder",
						Value:     "../src",
						TakesFile: true,
					},
				},
				Action: createSceneTableAction.CreateSceneTableAction,
			},
			{
				Name:  "update-scene-occurrence-subtitles",
				Usage: "Updates the character occurrences in scene's subtitles",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "src",
						Usage:     "Path to the latex source folder",
						Value:     "../src",
						TakesFile: true,
					},
				},
				Action: occurrenceSubtitles.UpdateCharacterOccurrenceSubtitles,
			},
			{
				Name:  "extract-page-numbers",
				Usage: "Create scene table",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:      "src",
						Usage:     "Path to the latex source folder",
						Value:     "../src",
						TakesFile: true,
					},
					&cli.StringFlag{
						Name:      "pdf",
						Usage:     "Path to current PDF file",
						Required:  true,
						TakesFile: true,
					},
				},
				Action: extractPageNumbers.ExtractPageNumbers,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
