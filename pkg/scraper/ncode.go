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
}

func NCodeInit(initial_param string) (*NCode, error) {
	if isUrl(initial_param) {
		var problem_slug = get_problem_slug(initial_param)
		return get_problem(problem_slug)
	}
	return nil, nil
}

func get_problem_slug(url string) string {
	vals := strings.Split(url, "/")
	return vals[len(vals)-1]
}

func get_problem(slug string) (*NCode, error) {
	_, err := check_cache(slug, makeRequest)
	if err != nil {
		println(err.Error())
	}
	return nil, nil
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
