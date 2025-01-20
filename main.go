package main

import (
	"fmt"
	"strings"
	"os"

	// "os/exec"
	"democli/start/utils"

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
	GROQ_API_KEY	string
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

	// if api not provided during build time
	if GROQ_API_KEY == "" {
		// load the env
		err := godotenv.Load("~/$HOME/Desktop/projects/gocli/.env")
		if err != nil {
			fmt.Println("Error loading .env file", err)
			return
		}
		GROQ_API_KEY = os.Getenv("GROQ_API_KEY")
	}

	// check if the user wats to chat or commit
	var purpose string = "commit"
	purForm := huh.NewSelect[string]().
		Title("Whats the mood?").
		Options(
			huh.NewOption("generate commit msg!", "commit"),
			huh.NewOption("anything else?", "chat"),
		).
		Value(&purpose)

	err := purForm.Run()
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

	analyzer := NewDiffAnalyzer()
	// its only for the chat purpose
	var user_query string
	if purpose == "chat" {
		inputForm := huh.NewInput().
			Title("Something different? ").
			Placeholder("Tell me a joke!").
			Value(&user_query)

		runErr := inputForm.Run()
		if runErr != nil {
			fmt.Println(runErr)
			return
		}
		var llm_context string
		context, err := SearchDirectory(".")
		if err != nil {
			fmt.Println(err)
		}
		
		for t := 0; t < len(context); t++ {
			llm_context = llm_context + context[t].FilePath + "\n"
			llm_context = llm_context + "Imports: " + strings.Join(context[t].Imports, ", ") + "\n"
			llm_context = llm_context + "Functions: " + strings.Join(context[t].Functions, ", ") + "\n"
			llm_context = llm_context + "Variables: " + strings.Join(context[t].Variables, ", ") + "\n"
		}

		Prompt = "Context: " + llm_context + "\n\nUser Query: " + user_query
	};if purpose == "commit" {
		cmd, err := utils.RunCommand("git", "diff")
		if err != nil {
			fmt.Println(err)
			return
		}
		
		Prompt, err = analyzer.analyzeGitDiff(cmd)
		if err != nil {
			fmt.Println(err)
			return
		}
	}


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
		// if purpose == "commit" {
		clipErr := clipboard.WriteAll(res)
		if clipErr != nil {
			fmt.Println("failed to copy ", clipErr)
			return
		}
		// }
	} else {
		fmt.Println("yare yare! Prompt is empty!!!")
	}
}
