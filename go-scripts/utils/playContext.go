package utils

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"regexp"
	"slices"
	"strconv"
)

type PlayContext struct {
	srcPath         string
	actTexFileNames []string
	Characters      []string
}

func GetPlayContext(srcPath string) (*PlayContext, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("couldn't find CWD: %w", err)
	}
	return collectFiles(path.Join(cwd, srcPath))
}

func (c PlayContext) GetSrcPath() string {
	return c.srcPath
}

func (c PlayContext) GetActTexFileNames() []string {
	return c.actTexFileNames
}

func collectFiles(srcPath string) (*PlayContext, error) {
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

	return &PlayContext{
			srcPath,
			actTexFileNames,
			characters,
		},
		nil
}

func parseCharacters(srcPath string) ([]string, error) {
	countersFilePath := path.Join(srcPath, "character_counters.tex")
	countersFileContent, err := os.ReadFile(countersFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't read character counters: %w", err)
	}

	characterRegex := regexp.MustCompile(`(?i)\\newcounter\{(?P<Character>[a-z]+)}`)

	var characters []string
	for _, match := range characterRegex.FindAllStringSubmatch(string(countersFileContent), -1) {
		characters = append(characters, match[1])
	}

	if len(characters) == 0 {
		return nil, fmt.Errorf("couldn't find any Characters in %s", countersFilePath)
	}

	slices.Sort(characters)

	log.Printf("Found Characters: %+v", characters)

	return characters, nil
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
