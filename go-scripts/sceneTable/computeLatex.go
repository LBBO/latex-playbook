package createSceneTableAction

import (
	"embed"
	"fmt"
	"github.com/lbbo/latex-playbook/go-scripts/utils"
	"strings"
	"text/template"
)

type sceneDescription struct {
	SceneShortName   string
	OccurrenceValues []string
}

type templateData struct {
	Characters  []string
	Scenes      []sceneDescription
	SceneCounts []uint
}

//go:embed templates/*
var templates embed.FS

func computeLatex(characters []string, actOccurrences []ActOccurrences) (string, error) {
	capitalizedCharacters, err := utils.CapitalizeCharacters(characters)

	if err != nil {
		return "", fmt.Errorf("could not capitalize characters: %w", err)
	}

	data := templateData{
		Characters:  capitalizedCharacters,
		Scenes:      computeSceneDescriptors(actOccurrences, characters),
		SceneCounts: countTotalScenesPerCharacter(characters, actOccurrences),
	}

	return renderLatex(data)
}

const latexForOccurrence = "\\cellcolor{TableColorAppearance}"

func computeSceneDescriptors(actOccurrences []ActOccurrences, characters []string) []sceneDescription {
	sceneDescriptors := make([]sceneDescription, 0)
	for actIndex, actOccurrence := range actOccurrences {
		for sceneIndex, scene := range actOccurrence {
			sceneShortName := fmt.Sprintf("%d.%d", actIndex+1, sceneIndex+1)
			sceneShortName = fmt.Sprintf("%5s", sceneShortName)
			occurrenceValues := make([]string, 0)

			for _, character := range characters {
				occurrenceValue := strings.Repeat(" ", len(latexForOccurrence))

				if scene[strings.ToLower(character)] > 0 {
					occurrenceValue = fmt.Sprintf("%s %d", latexForOccurrence, scene[strings.ToLower(character)])
				}

				occurrenceValues = append(occurrenceValues, occurrenceValue)
			}

			sceneDescriptors = append(sceneDescriptors, sceneDescription{
				SceneShortName:   sceneShortName,
				OccurrenceValues: occurrenceValues,
			})
		}
	}
	return sceneDescriptors
}

func renderLatex(data templateData) (string, error) {
	file, err := templates.ReadFile("templates/tableTex.gotxt")
	if err != nil {
		return "", fmt.Errorf("could not read template file: %w", err)
	}
	templ, err := template.New("sceneTable.tex").Parse(string(file))
	if err != nil {
		return "", fmt.Errorf("could not parse template file: %w", err)
	}

	result := new(strings.Builder)
	err = templ.Execute(result, data)
	if err != nil {
		return "", fmt.Errorf("could not execute template: %w", err)
	}
	return result.String(), nil
}

func countTotalScenesPerCharacter(characters []string, actOccurrences []ActOccurrences) []uint {
	sceneCounts := make([]uint, len(characters))
	for _, actOccurrence := range actOccurrences {
		for _, scene := range actOccurrence {
			for characterIndex, character := range characters {
				if scene[strings.ToLower(character)] > 0 {
					sceneCounts[characterIndex]++
				}
			}
		}
	}
	return sceneCounts
}
