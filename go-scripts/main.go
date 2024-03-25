package main

import (
	"log"
	"os"

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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
