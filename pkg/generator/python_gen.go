package generator

import (
	"bytes"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/BFostek/ProvaGen/pkg/scraper"
	"github.com/BFostek/ProvaGen/pkg/utils"
)

type PythonGenerator struct {
	scrapper scraper.Scrapper
}

func (pg *PythonGenerator) Generate(project_destination, challenge_id string) error {
	var err error
	pg.scrapper, err = scraper.NCodeInit(challenge_id)
	if err != nil {
		return err
	}
	var val scraper.Challenge
	if val, err = pg.scrapper.GetChallenge(); err != nil {
		return err
	}
	if err := utils.CreateStructure(project_destination, challenge_id); err != nil {
		return err
	}
	destination := path.Join(project_destination, challenge_id)
	pg.CreateFiles(val, destination)
	println(val.Name)

	return nil
}

func (pg *PythonGenerator) CreateFiles(val scraper.Challenge, destination string) {
	if val.InitialFile != nil {
		mainContent := pg.includeMainContent(*val.InitialFile)
		fullContent := *val.InitialFile + mainContent
		utils.CreateFileWithContent(path.Join(destination, "main.py"), fullContent)
	}
	if val.Solution != nil {
		utils.CreateFileWithContent(path.Join(destination, "solution.py"), *val.Solution)
	}
	if val.Description != nil {
		//TODO Parse description file
		utils.CreateFileWithContent(path.Join(destination, "description"), *val.Description)
	}
	if val.Tests != nil {
		for idx, item := range val.Tests {
			for l, v := range item {
				println(idx, l, v)
			}

		}
		utils.CreateFileWithContent(path.Join(destination, "main_test.py"), "")
		//TODO test
	}
}
func (pg *PythonGenerator) includeMainContent(s string) string {
	mainBlock := `
if __name__ == "__main__":
    sol = Solution()
    # Test case structure (implement or add your own)
    # Method names found: {{ join .Methods ", " }}
    # Example usage:
    # result = sol.{{ .FirstMethod }}(*args)
    # print(result)`

	// Create a template function map with the join function
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	// Parse initial code to find Solution methods
	methods := pg.extractSolutionMethods(s)

	// Create template data
	data := struct {
		Methods     []string
		FirstMethod string
	}{
		Methods:     methods,
		FirstMethod: pg.firstOrDefault(methods),
	}

	// Create template with function map
	tpl := template.New("main").Funcs(funcMap)
	tpl = template.Must(tpl.Parse(mainBlock))

	var mainContent bytes.Buffer
	if err := tpl.Execute(&mainContent, data); err != nil {
		println("Template execution error: %v", err)
		return ""
	}
	return mainContent.String()
}

// Helper to extract Solution class method names
func (pg *PythonGenerator) extractSolutionMethods(code string) []string {
	var methods []string
	re := regexp.MustCompile(`def\s+(\w+)\s*\(`)
	matches := re.FindAllStringSubmatch(code, -1)

	for _, match := range matches {
		if len(match) > 1 {
			methods = append(methods, match[1])
		}
	}
	return methods
}

// Helper to get first method or placeholder
func (pg *PythonGenerator) firstOrDefault(methods []string) string {
	if len(methods) > 0 {
		return methods[0]
	}
	return "methodName"
}
