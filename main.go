package main

import (
	"fmt"
	"strings"

	// "os/exec"
	"democli/start/utils"
	"os"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	Model           string
	Prompt          string
	AvailableModels []string
	ModelsMapped    map[string]interface{}
	modelOptions    []huh.Option[string]
)

func loadHuhModelOption(API_KEY string) {
	AvailableModels = utils.GetAvailableModels(API_KEY)
	// huh options
	for _, model := range AvailableModels {
		modelName := cases.Title(language.English, cases.NoLower).String(strings.Join(strings.Split(model, "-"), " "))
		modelOptions = append(modelOptions, huh.NewOption(modelName, model))
	}
}

func main() {

	// load the env
	err := godotenv.Load("/home/shrek/Desktop/projects/gocli/.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
	GROQ_API_KEY := os.Getenv("GROQ_API_KEY")

	// check if the user wats to chat or commit
	var purpose string
	purForm := huh.NewSelect[string]().
		Title("Whats the mood?").
		Options(
			huh.NewOption("chat", "chat"),
			huh.NewOption("generate commit msg!", "commit"),
		).
		Value(&purpose)

	err = purForm.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	// fetch the available models
	loadHuhModelOption(GROQ_API_KEY)

	selectModelForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a Model: ").
				Options(
					modelOptions...,
				).
				Value(&Model),
		),
	)

	runErr := selectModelForm.Run()
	if runErr != nil {
		fmt.Println(runErr)
		return
	}
	// its onlly for the chat purpose
	if purpose == "chat" {
		inputForm := huh.NewInput().
			Title("Something different? ").
			Placeholder("Tell me a joke!").
			Value(&Prompt)

		runErr := inputForm.Run()
		if runErr != nil {
			fmt.Println(runErr)
			return
		}
	}

	Prompt = getPrompt(purpose, Prompt)

	// fmt.Println(Prompt)
	if Prompt != "" {
		res := utils.ModelCall(Model, Prompt, "", GROQ_API_KEY)
		if res == "" {
			fmt.Println("Model Calling went wrong")
			return
		}
		if res[0] == '"' {
			res = res[1 : len(res)-1]
		}
		fmt.Println(res)
	
		// copy to clipboard
		if purpose == "commit" {
			clipErr := clipboard.WriteAll(res)
			if clipErr != nil {
				fmt.Println("failed to copy ", clipErr)
				return
			}
		}
	} else {
		fmt.Println("yare yare! Prompt is empty!!!")
	}
}
