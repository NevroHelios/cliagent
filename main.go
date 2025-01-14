package main

import (
	"fmt"
	"strings"

	// "os/exec"
	"democli/start/utils"
	"os"
	"regexp"

	"github.com/WqyJh/go-fstring"
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
)

func main() {

	// load the env
	err := godotenv.Load("/home/shrek/Desktop/projects/gocli/.env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
	}
	GROQ_API_KEY := os.Getenv("GROQ_API_KEY")
	// get available models
	get_models_url := "https://api.groq.com/openai/v1/models"
	AvailableModels = utils.GetAvailableModels(get_models_url, GROQ_API_KEY)

	// huh options
	modelOptions := []huh.Option[string]{}
	for _, model := range AvailableModels {
		modelName := cases.Title(language.English, cases.NoLower).String(strings.Join(strings.Split(model, "-"), " "))
		modelOptions = append(modelOptions, huh.NewOption(modelName, model))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a Model: ").
				Options(
					modelOptions...,
				).
				Value(&Model),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Something different? ").
				Placeholder("Tell me a joke!").
				Value(&Prompt),
		),
	)

	runErr := form.Run()
	if runErr != nil {
		fmt.Println(runErr)
		return
	}

	git_diff := ""
	if Prompt == "" {
		// get the git diff
		cmd, err := utils.RunCommand("git", "diff")
		if err != nil {
			fmt.Println("No git diff", err)
			return
		}

		template := `
		You are an expert software engineer assisting with writing clear and concise git commit messages. Given the following "git diff" output, analyze the changes and provide a descriptive commit message summarizing the purpose and impact of the modifications in 50 words or less. 
		NOTE: You are to return only the commit message and nothing else.
		
		Git Diff:
		{git_diff}
		
		Commit Message:
		
		`
		lines := strings.Split(cmd, "\n")
		if len(lines) > 1000 {
			re := regexp.MustCompile(`^.func`)
			for _, line := range lines {
				if re.MatchString(line) {
					git_diff += line + "\n"
				}
			}
		} else {
			git_diff = cmd
		}

		diff := map[string]any{"git_diff": git_diff}
		Prompt, err = fstring.Format(template, diff)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// fmt.Println(Prompt)
	res := utils.ModelCall(Model, Prompt, "", GROQ_API_KEY)
	if res == "" {
		fmt.Println("Something went wrong")
		return
	}
	if res[0] == '"' {
		res = res[1 : len(res)-1]
	}
	fmt.Println(res)

	// copy to clipboard
	if git_diff != "" {
		clipErr := clipboard.WriteAll(res)
		if clipErr != nil {
			fmt.Println("Error initializing clipboard:", clipErr)
			return
		}
	}
}
