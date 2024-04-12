package occurrenceSubtitles

import (
	"fmt"
	createSceneTableAction "github.com/lbbo/latex-playbook/go-scripts/sceneTable"
	"github.com/lbbo/latex-playbook/go-scripts/utils"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"regexp"
	"strings"
)

func UpdateCharacterOccurrenceSubtitles(c *cli.Context) error {
	playContext, err := utils.GetPlayContext(c.String("src"))
	if err != nil {
		return err
	}

	actOccurrences, err := createSceneTableAction.ComputeOccurrences(playContext)
	if err != nil {
		return fmt.Errorf("couldn't compute occurrences: %w", err)
	}

	capitalizedCharacters, err := utils.CapitalizeCharacters(playContext.Characters)
	if err != nil {
		return fmt.Errorf("couldn't capitalize characters: %w", err)
	}

	for actIndex, actTexFileName := range playContext.GetActTexFileNames() {
		actTexFilePath := path.Join(playContext.GetSrcPath(), actTexFileName)
		err := updateAct(actTexFilePath, capitalizedCharacters, actOccurrences[actIndex])
		if err != nil {
			return fmt.Errorf("couldn't update act %d: %w", actIndex, err)
		}
	}

	return nil
}

func updateAct(actTexFilePath string, characters []string, actOccurrences createSceneTableAction.ActOccurrences) error {
	actTexFile, err := os.ReadFile(actTexFilePath)
	if err != nil {
		return fmt.Errorf("couldn't read act file: %w", err)
	}

	fileContent := string(actTexFile)
	rest := fileContent
	for _, scene := range actOccurrences {
		var available []string
		for _, character := range characters {
			if scene[strings.ToLower(character)] > 0 {
				available = append(available, fmt.Sprintf("%s (%d)", character, scene[strings.ToLower(character)]))
			}
		}
		newCharactersList := strings.Join(available, ", ")

		rest = rest[:]
		regex := regexp.MustCompile(`\\scene\{(.*?)}\{(.*?)}`)
		matchIndexes := regex.FindStringSubmatchIndex(rest)
		if matchIndexes != nil {
			sceneTitleTex := strings.Replace(rest[matchIndexes[0]:matchIndexes[1]], rest[matchIndexes[4]:matchIndexes[5]], newCharactersList, 1)
			fileContent = strings.Replace(fileContent, rest[matchIndexes[0]:matchIndexes[1]], sceneTitleTex, 1)
			rest = rest[matchIndexes[1]:]
		}
	}

	err = os.WriteFile(actTexFilePath, []byte(fileContent), 0644)

	if err != nil {
		return fmt.Errorf("couldn't write file: %w", err)
	}

	return nil
}
