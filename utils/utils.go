package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetAvailableModels(url string, API_KEY string) []string {
	var availableModels []string 

	req, err := http.NewRequest("GET", url, nil)
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

	return availableModels
}
