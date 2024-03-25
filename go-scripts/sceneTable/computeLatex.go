package createSceneTableAction

import (
	"embed"
	"fmt"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

type sceneDescription struct {
	SceneShortName   string
	OccurrenceValues []string
}

type templateData struct {
	Characters []string
	Scenes     []sceneDescription
}

//go:embed templates/*
var templates embed.FS

func computeLatex(characters []string, actOccurrences []actOccurrences) (string, error) {
	var capitalizedCharacters []string

	for _, character := range characters {
		r, size := utf8.DecodeRuneInString(character)
		if r == utf8.RuneError {
			return "", fmt.Errorf("could not decode rune in string: %w", r)
		}
		capitalizedCharacters = append(capitalizedCharacters, string(unicode.ToUpper(r))+character[size:])
	}

	data := templateData{
		Characters: capitalizedCharacters,
		Scenes:     computeSceneDescriptors(actOccurrences, characters),
	}

	return renderLatex(data)
}

const latexForOccurrence = "\\cellcolor{TableColorAppearance}"

func computeSceneDescriptors(actOccurrences []actOccurrences, characters []string) []sceneDescription {
	sceneDescriptors := make([]sceneDescription, 0)
	for actIndex, actOccurrence := range actOccurrences {
		for sceneIndex, scene := range actOccurrence {
			sceneShortName := fmt.Sprintf("%d.%d", actIndex+1, sceneIndex+1)
			sceneShortName = fmt.Sprintf("%5s", sceneShortName)
			occurrenceValues := make([]string, 0)

			for _, character := range characters {
				occurrenceValue := strings.Repeat(" ", len(latexForOccurrence))

				if scene[strings.ToLower(character)] {
					occurrenceValue = "\\cellcolor{TableColorAppearance}"
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
