package createSceneTableAction

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
)

type context struct {
	srcPath         string
	actTexFileNames []string
	characters      []string
}

type actOccurrences []map[string]bool

func CreateSceneTableAction(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("couldn't find CWD: %w", err)
	}
	context, err := collectFiles(path.Join(cwd, c.String("src")))

	if err != nil {
		return err
	}

	actOccurrences, err := context.computeOccurrences()

	if err != nil {
		return err
	}

	log.Printf("all occurrences: %+v", actOccurrences)

	latex, err := computeLatex(context.characters, actOccurrences)
	if err != nil {
		return err
	}

	return context.writeTable(latex)
}

func (c *context) computeOccurrences() ([]actOccurrences, error) {
	var result []actOccurrences
	for _, actTexFileName := range c.actTexFileNames {
		actTexFilePath := path.Join(c.srcPath, actTexFileName)
		var actTexFile io.Reader
		actTexFile, err := os.OpenFile(actTexFilePath, os.O_RDONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("couldn't read act file: %w", err)
		}
		reader := bufio.NewScanner(actTexFile)

		occurrences, err := computeOccurrences(reader)
		if err != nil {
			return nil, fmt.Errorf("couldn't compute occurrences: %w", err)
		}
		result = append(result, occurrences)
	}
	return result, nil
}

func (c *context) writeTable(tex string) error {
	filePath := path.Join(c.srcPath, "occurrence_page.tex")
	err := os.WriteFile(filePath, []byte(tex), 0644)

	if err != nil {
		return fmt.Errorf("couldn't update occurrence page: %w", err)
	}

	return nil
}

func computeOccurrences(reader *bufio.Scanner) (actOccurrences, error) {
	characterRegex := regexp.MustCompile("(?i)\\\\characters*\\{(?P<Characters>(?:[a-z]|}\\{)+)}")
	sceneRegex := regexp.MustCompile("(?i)\\\\scene")
	var scenes actOccurrences
	var currScene *map[string]bool
	for reader.Scan() {
		line := reader.Text()
		if sceneRegex.MatchString(line) {
			scenes = append(scenes, make(map[string]bool))
			currScene = &scenes[len(scenes)-1]
		} else if characterRegex.MatchString(line) {
			characters := strings.Split(characterRegex.FindStringSubmatch(line)[1], "}{")

			if currScene == nil {
				return nil, fmt.Errorf("characters %q appeared before first scene", strings.Join(characters, ", "))
			}

			for _, character := range characters {
				(*currScene)[character] = true
			}
		}
	}

	return scenes, nil
}

func parseCharacters(srcPath string) ([]string, error) {
	countersFilePath := path.Join(srcPath, "character_counters.tex")
	countersFileContent, err := os.ReadFile(countersFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read character counters: %w", err)
	}

	characterRegex := regexp.MustCompile("(?i)\\\\newcounter\\{(?P<Character>[a-z]+)}")

	var characters []string
	for _, match := range characterRegex.FindAllStringSubmatch(string(countersFileContent), -1) {
		characters = append(characters, match[1])
	}

	if len(characters) == 0 {
		return nil, fmt.Errorf("couldn't find any characters in %s", countersFilePath)
	}

	slices.Sort(characters)

	log.Printf("Found characters: %+v", characters)

	return characters, nil
}

func collectFiles(srcPath string) (*context, error) {
	srcEntries, err := os.ReadDir(srcPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read latex Floh im Ohr folder: %w", err)
	}

	characters, err := parseCharacters(srcPath)

	if err != nil {
		return nil, err
	}

	actTexFileNames, err := getSortedActFileNames(srcEntries)
	if err != nil {
		return nil, err
	}
	log.Printf("Will parse these act files: %+v", actTexFileNames)

	return &context{
			srcPath,
			actTexFileNames,
			characters,
		},
		nil
}

func getSortedActFileNames(srcEntries []os.DirEntry) ([]string, error) {
	fileNamesMap := make(map[int]string)
	minIndex := math.MaxInt
	maxIndex := math.MinInt
	actRegex := regexp.MustCompile("(?i)act_(?P<Index>[0-9]+).tex")
	for _, entry := range srcEntries {
		if actRegex.MatchString(entry.Name()) && !entry.IsDir() {
			index, err := strconv.Atoi(actRegex.FindStringSubmatch(entry.Name())[1])
			if err != nil {
				return nil, fmt.Errorf("couldn't parse act index: %w", err)
			}
			fileNamesMap[index] = entry.Name()
			if index < minIndex {
				minIndex = index
			}
			if index > maxIndex {
				maxIndex = index
			}
		}
	}

	if minIndex != 1 {
		return nil, fmt.Errorf("act files should start from 1, but found %d", minIndex)
	}

	var actTexFileNames []string
	for i := minIndex; i <= maxIndex; i++ {
		if _, ok := fileNamesMap[i]; !ok {
			return nil, fmt.Errorf("missing act file for index %d", i)
		}
		actTexFileNames = append(actTexFileNames, fileNamesMap[i])
	}

	return actTexFileNames, nil
}
