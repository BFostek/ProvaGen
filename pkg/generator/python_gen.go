package generator

import (
	"bytes"
	"fmt"
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
	imports_packs := `
from collections import *
from itertools import *
from typing import List, Dict, Set, Optional, Protocol
from copy import deepcopy
import traceback
import heapq
import collections
import sys
import re
import math
import ast
`

	if val.InitialFile != nil {
		mainContent := pg.includeMainContent(*val.InitialFile)
		fullContent := imports_packs + *val.InitialFile + mainContent
		utils.CreateFileWithContent(path.Join(destination, "main.py"), fullContent)
	}
	if val.Solution != nil {
		utils.CreateFileWithContent(path.Join(destination, "solution.py"), imports_packs+*val.Solution)
	}
	if val.Description != nil {
		//TODO Parse description file
		utils.CreateFileWithContent(path.Join(destination, "description"), *val.Description)
	}
	if val.Tests != nil {
		testContent := "import pytest\n"
		testContent += "from main import Solution as SolutionMain\n"
		testContent += "from solution import Solution as SolutionSolution\n\n"

		// Extract method name from solution code
		solutionCode := *val.Solution
		methods := pg.extractSolutionMethods(solutionCode)
		if len(methods) == 0 {
			return // No methods found
		}
		methodName := methods[0]

		// Extract parameters for the method (excluding self)
		params := pg.extractMethodParams(solutionCode, methodName)

		for idx, item := range val.Tests {
			// Collect arguments in parameter order
			args := make([]string, 0, len(params))
			for _, param := range params {
				value, exists := item[param]
				if !exists {
					fmt.Printf("Missing param '%s' in test case %d\n", param, idx)
					break
				}
				args = append(args, pg.convertValueToPython(value))
			}

			if len(args) != len(params) {
				continue // Skip incomplete test cases
			}

			testContent += fmt.Sprintf("def test_case_%d():\n", idx)
			for i, param := range params {
				testContent += fmt.Sprintf("    %s = %s\n", param, args[i])
			}
			testContent += "    main_sol = SolutionMain()\n"
			testContent += "    solution_sol = SolutionSolution()\n"
			argsCall := strings.Join(params, ", ")
			testContent += fmt.Sprintf(
				"    assert main_sol.%s(%s) == solution_sol.%s(%s)\n\n",
				methodName, argsCall, methodName, argsCall,
			)

		}
		utils.CreateFileWithContent(path.Join(destination, "main_test.py"), testContent)
		//TODO test
	}
}
func (pg *PythonGenerator) extractMethodParams(code, methodName string) []string {
	re := regexp.MustCompile(`def\s+` + regexp.QuoteMeta(methodName) + `\s*\(([^)]*)\)`)
	matches := re.FindStringSubmatch(code)
	if len(matches) < 2 {
		return nil
	}

	params := strings.Split(matches[1], ",")
	paramsList := make([]string, 0)
	for _, p := range params {
		p = strings.TrimSpace(p)
		if p == "self" || p == "" {
			continue
		}
		// Remove type hints/default values (e.g., "intervals: List[List[int]]" -> "intervals")
		name := strings.SplitN(p, ":", 2)[0]
		name = strings.SplitN(name, "=", 2)[0]
		paramsList = append(paramsList, strings.TrimSpace(name))
	}
	return paramsList
}

// Converts Go values to Python-compatible strings (e.g., slices to lists)
func (pg *PythonGenerator) convertValueToPython(value any) string {
	switch v := value.(type) {
	case []any:
		elements := make([]string, len(v))
		for i, e := range v {
			elements[i] = pg.convertValueToPython(e)
		}
		return "[" + strings.Join(elements, ", ") + "]"
	case string:
		return fmt.Sprintf("%q", v)
	default:
		return fmt.Sprintf("%v", v)
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
