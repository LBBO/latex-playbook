package main

import (
	"log"
	"os"

	"github.com/lbbo/latex-playbook/go-scripts/extractPageNumbers"
	occurrenceSubtitles "github.com/lbbo/latex-playbook/go-scripts/occurrence-subtitles"

	createSceneTableAction "github.com/lbbo/latex-playbook/go-scripts/sceneTable"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "latex playbook scripts",
		Usage: "Scripts for managing latex playbook",
		Commands: []*cli.Command{
			{
				Name:  "create-scene-table",
				Usage: "Create scene table",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:  "src",
						Usage: "Path to the latex source folder",
						Value: "../src",
					},
				},
				Action: createSceneTableAction.CreateSceneTableAction,
			},
			{
				Name:  "update-scene-occurrence-subtitles",
				Usage: "Updates the character occurrences in scene's subtitles",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:  "src",
						Usage: "Path to the latex source folder",
						Value: "../src",
					},
				},
				Action: occurrenceSubtitles.UpdateCharacterOccurrenceSubtitles,
			},
			{
				Name:  "extract-page-numbers",
				Usage: "Create scene table",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:  "src",
						Usage: "Path to the latex source folder",
						Value: "../src",
					},
					&cli.PathFlag{
						Name:     "pdf",
						Usage:    "Path to current PDF file",
						Required: true,
					},
				},
				Action: extractPageNumbers.ExtractPageNumbers,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
