package createSceneTableAction

import (
	"bufio"
	"context"
	"fmt"
	"github.com/lbbo/latex-playbook/go-scripts/utils"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/urfave/cli/v3"
)

type ActOccurrences []map[string]uint

type FileNamesGetter interface {
	GetActTexFileNames() []string
	GetSrcPath() string
}

func CreateSceneTableAction(ctx context.Context, c *cli.Command) error {
	playContext, err := utils.GetPlayContext(c.String("src"))

	if err != nil {
		return err
	}

	actOccurrences, err := ComputeOccurrences(playContext)

	if err != nil {
		return err
	}

	log.Printf("all occurrences: %+v", actOccurrences)

	latex, err := computeLatex(playContext.Characters, actOccurrences)
	if err != nil {
		return err
	}

	return writeTable(playContext, latex)
}

func ComputeOccurrences(c FileNamesGetter) ([]ActOccurrences, error) {
	var result []ActOccurrences
	for _, actTexFileName := range c.GetActTexFileNames() {
		actTexFilePath := path.Join(c.GetSrcPath(), actTexFileName)
		var actTexFile io.Reader
		actTexFile, err := os.OpenFile(actTexFilePath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("couldn't read act file: %w", err)
		}
		reader := bufio.NewScanner(actTexFile)

		occurrences, err := computeOccurrencesFromFile(reader)
		if err != nil {
			return nil, fmt.Errorf("couldn't compute occurrences: %w", err)
		}
		result = append(result, occurrences)
	}
	return result, nil
}

func writeTable(c FileNamesGetter, tex string) error {
	filePath := path.Join(c.GetSrcPath(), "occurrence_page.tex")
	err := os.WriteFile(filePath, []byte(tex), 0644)

	if err != nil {
		return fmt.Errorf("couldn't update occurrence page: %w", err)
	}

	return nil
}

func computeOccurrencesFromFile(reader *bufio.Scanner) (ActOccurrences, error) {
	characterRegex := regexp.MustCompile(`(?i)\\characters*\{(?P<Characters>(?:[a-z]|}\{)+)}`)
	sceneRegex := regexp.MustCompile(`(?i)\\scene`)
	var scenes ActOccurrences
	var currScene *map[string]uint
	for reader.Scan() {
		line := reader.Text()
		if sceneRegex.MatchString(line) {
			scenes = append(scenes, make(map[string]uint))
			currScene = &scenes[len(scenes)-1]
		} else if characterRegex.MatchString(line) {
			characters := strings.Split(characterRegex.FindStringSubmatch(line)[1], "}{")

			if currScene == nil {
				return nil, fmt.Errorf("characters %q appeared before first scene", strings.Join(characters, ", "))
			}

			for _, character := range characters {
				(*currScene)[character]++
			}
		}
	}

	return scenes, nil
}
