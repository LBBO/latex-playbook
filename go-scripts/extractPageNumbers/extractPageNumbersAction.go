package extractPageNumbers

import (
	"fmt"
	"github.com/lbbo/latex-playbook/go-scripts/utils"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/urfave/cli/v2"
	"os"
	"path"
	"regexp"
)

func ExtractPageNumbers(c *cli.Context) error {
	srcPath := c.Path("src")
	pdfPath := c.Path("pdf")

	allPageNumbers, err := getPageNumbers(srcPath, pdfPath)
	if err != nil {
		return fmt.Errorf("couldn't get page numbers: %w", err)
	}

	playContext, err := utils.GetPlayContext(srcPath)
	if err != nil {
		return fmt.Errorf("couldn't get play context: %w", err)
	}

	filePath := path.Join(playContext.GetSrcPath(), "occurrence_page.tex")
	buff, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("couldn't read file: %w", err)
	}
	fileContent := string(buff)

	for actIndex, actPageNumbers := range allPageNumbers {
		for sceneIndex, pageNumbers := range actPageNumbers {
			regex := regexp.MustCompile(fmt.Sprintf(`(%d\.%d\s*)&(.*?)&`, actIndex+1, sceneIndex+1))
			fileContent = regex.ReplaceAllString(fileContent, fmt.Sprintf("$1& %d - %d &", pageNumbers[0], pageNumbers[1]))
		}
	}

	err = os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return fmt.Errorf("couldn't write file: %w", err)
	}

	return nil
}

func getPageNumbers(srcPath string, pdfPath string) ([][][]int, error) {
	if srcPath == "" {
		return nil, fmt.Errorf("src flag is required")
	}

	if pdfPath == "" {
		return nil, fmt.Errorf("pdf flag is required")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("couldn't find CWD: %w", err)
	}

	completePdfPath := path.Join(cwd, pdfPath)
	fmt.Println(completePdfPath)

	pdf, err := os.Open(completePdfPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open pdf file: %w", err)
	}
	defer pdf.Close()

	pdfInfo, err := api.PDFInfo(pdf, path.Base(completePdfPath), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't get pdf info: %w", err)
	}

	bookmarks, err := api.Bookmarks(pdf, nil)
	if err != nil {
		return nil, fmt.Errorf("couldn't get pdf info: %w", err)
	}

	acts := bookmarks[1:]
	// -1 because we want 1-based page numbers!
	firstPageOffset := acts[0].PageFrom - 1
	var actPages [][][]int
	for actIndex, act := range acts {
		scenePages := make([][]int, 0)

		// last page number will be 0, so override with PDF's last page
		lastActPage := act.PageThru
		if actIndex == len(acts)-1 {
			lastActPage = pdfInfo.PageCount
		}

		for sceneIndex, scene := range act.Kids {
			// last page number will be 0, so override with act's last page
			lastScenePage := scene.PageThru + 1
			if sceneIndex == len(act.Kids)-1 {
				lastScenePage = lastActPage
			}
			scenePages = append(scenePages, []int{scene.PageFrom - firstPageOffset, lastScenePage - firstPageOffset})
		}

		actPages = append(actPages, scenePages)
	}
	return actPages, nil
}
