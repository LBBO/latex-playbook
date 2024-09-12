package occurrenceSubtitles

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	createSceneTableAction "github.com/lbbo/latex-playbook/go-scripts/sceneTable"
	"github.com/lbbo/latex-playbook/go-scripts/utils"
	"github.com/urfave/cli/v2"
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
		// The `}{` is moved into the second capture group to allow the second `{}`
		// to be empty and still be able to find the capture group uniquely within
		// the entire match so that we can insert the new character list into it.
		regex := regexp.MustCompile(`\\scene\{(.*?)(}\{.*?)}`)
		matchIndexes := regex.FindStringSubmatchIndex(rest)
		if matchIndexes != nil {
			oldSceneTitleTex := rest[matchIndexes[0]:matchIndexes[1]]
			oldCharacterList := rest[matchIndexes[4]:matchIndexes[5]]
			newSceneTitleTex := strings.Replace(
				oldSceneTitleTex,
				oldCharacterList,
				// we have to re-insert the }{ that was "removed" by the regex
				fmt.Sprintf("}{%s", newCharactersList),
				1,
			)
			fileContent = strings.Replace(
				fileContent,
				oldSceneTitleTex,
				newSceneTitleTex,
				1,
			)
			rest = rest[matchIndexes[1]:]
		}
	}

	err = os.WriteFile(actTexFilePath, []byte(fileContent), 0644)

	if err != nil {
		return fmt.Errorf("couldn't write file: %w", err)
	}

	return nil
}
