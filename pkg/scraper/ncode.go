package scraper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/BFostek/ProvaGen/configs"
)

type NCode struct {
	json_str string
}

type RawChallenge struct {
	Result struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Solutions   struct {
			Python string `json:"python"`
		} `json:"solutions"`
		InitialCode struct {
			Python string `json:"python"`
		} `json:"starterCode"`
		CustomTestCases []string `json:"custom_test_cases"`
	} `json:"result"`
}

func convertViaIntermediate(jsonData string) (*Challenge, error) {
	var raw RawChallenge
	if err := json.Unmarshal([]byte(jsonData), &raw); err != nil {
		return nil, err
	}

	challenge := &Challenge{
		Name: raw.Result.Name,
	}

	if raw.Result.Description != "" {
		challenge.Description = &raw.Result.Description
	}

	if raw.Result.InitialCode.Python != "" {
		println("Entrou aqui")
		challenge.InitialFile = &raw.Result.InitialCode.Python
	}
	if raw.Result.Solutions.Python != "" {
		challenge.Solution = &raw.Result.Solutions.Python
	}
	challenge.Tests = make([]map[string]string, len(raw.Result.CustomTestCases))
	for index, tc := range raw.Result.CustomTestCases {
		question_map := make(map[string]string)
		test_line := strings.Split(tc, "\n")
		for _, line := range test_line {
			l := strings.Split(line, "=")
			question_map[l[0]] = l[1]
		}
		challenge.Tests[index] = question_map
	}

	println(challenge.Name)
	return challenge, nil
}

func (ncode *NCode) GetChallenge() (Challenge, error) {
	var err error
	if val, err := convertViaIntermediate(ncode.json_str); err == nil {
		return *val, nil
	}

	return Challenge{}, err
}

func NCodeInit(initial_param string) (*NCode, error) {
	result := NCode{}
	var problem_slug = initial_param
	if isUrl(initial_param) {
		problem_slug = get_problem_slug(initial_param)
	}
	if problem, err := get_problem(problem_slug); err == nil {
		result.json_str = problem
		return &result, nil
	}
	return nil, nil
}

func get_problem_slug(url string) string {
	vals := strings.Split(url, "/")
	return vals[len(vals)-1]
}

func get_problem(slug string) (string, error) {
	return check_cache(slug, makeRequest)
}

func makeRequest(problemId string) (*string, error) {
	// Create request payload
	config, err := configs.LoadConfig("config.yaml")
	if err != nil {
		println(err.Error())
		return nil, err
	}
	payload := map[string]any{
		"data": map[string]string{
			"problemId": problemId,
		},
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest("POST",
		config.APIEndpoint.URL,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", config.APIEndpoint.Origin)
	req.Header.Set("Referer", config.APIEndpoint.Origin)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36")
	response, err := http.DefaultClient.Do(req)
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	result := string(body)
	return &result, nil
}

func isUrl(url string) bool {
	pattern := regexp.MustCompile(`^(https?://)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}/problems/[a-zA-Z0-9_-]+/?$`)
	return pattern.MatchString(url)
}
