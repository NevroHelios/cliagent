package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"democli/start/utils"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

var (
	url             string
	Model           string
	Messages        []Message
	AvailableModels []string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
	// fmt.Println(api_key)
	url = "https://api.groq.com/openai/v1/models"
	AvailableModels = utils.GetAvailableModels(url, os.Getenv("GROQ_API_KEY"))

	// for i := 0; i < len(AvailableModels); i++ {
	// 	fmt.Println(AvailableModels[i])
	// }
	fmt.Println(AvailableModels)
}
