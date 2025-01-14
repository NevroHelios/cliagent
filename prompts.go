package main

import (
	"democli/start/utils"
	"fmt"
	"regexp"
	"strings"

	"github.com/WqyJh/go-fstring"
)

var (
	git_diff string
)

func getPrompt(purpose string, prompt string) string {
	if purpose == "commit" {
		// get the git diff
		cmd, err := utils.RunCommand("git", "diff")
		if err != nil {
			fmt.Println("No git diff", err)
			return ""
		}

		template := `
				You are an expert software engineer assisting with writing clear and concise git commit messages. Given the following "git diff" output, analyze the changes and provide a descriptive commit message summarizing the purpose and impact of the modifications in 50 words or less. 
				NOTE: You are to return only the commit message and nothing else.
				
				Git Diff:
				{git_diff}
				
				Commit Message:
				
				`
		lines := strings.Split(cmd, "\n")
		if len(lines) > 2000 {
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
			fmt.Println("Error formatting prompt", err)
			return ""
		}
		return Prompt
	} else {
		template := `
			You are a helpful assistant that can answer questions about code.
			You should answer the question as concisely as possible.
			If you don't know the answer, just say that you don't know, don't try to make up an answer.
			Note: You are to return only the answer and nothing else.
			Question:
			{prompt}
			
			Answer:
			`
		dctPrompt := map[string]any{"prompt": prompt}
		Prompt, err := fstring.Format(template, dctPrompt)
		if err != nil {
			fmt.Println("Error formatting prompt", err)
			return ""
		}
		return Prompt
	}
}
