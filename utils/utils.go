package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"io"
	"net/http"
	"os/exec"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Body struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func GetAvailableModels(API_KEY string) []string {
	var availableModels []string
	get_models_url := "https://api.groq.com/openai/v1/models"

	req, err := http.NewRequest("GET", get_models_url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Add("Authorization", "Bearer "+API_KEY)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	if data, ok := response["data"].([]interface{}); ok {
		for _, item := range data {
			if model, ok := item.(map[string]interface{})["id"].(string); ok {
				availableModels = append(availableModels, model)
			}
		}
	} else {
		fmt.Println("No models available or incorrect response format.")
	}

	sort.Slice(availableModels, func(i, j int) bool {
		return availableModels[i] < availableModels[j]
	})

	return availableModels
}

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func ModelCall(model string, prompt string, sys_prompt string, API_KEY string) string {
	call_url := "https://api.groq.com/openai/v1/chat/completions"

	body := Body{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}
	req, err := http.NewRequest("POST", call_url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return ""
	}
	// headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+API_KEY)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return ""
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return ""
	}

	if choices, ok := response["choices"].([]interface{}); ok {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := choice["message"].(map[string]interface{})["content"].(string); ok {
				return text
			}
		}
	}
	return ""
}

