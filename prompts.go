package main

// import (
// 	"fmt"
// 	"regexp"
// 	"strings"
// )
import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ChangeType struct {
	Pattern     string
	Description string
	Priority    int
}

var languagePatterns = map[string][]ChangeType{
	"go": {
		{Pattern: `^[+-]\s*func\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*type\s+\w+`, Description: "type definitions", Priority: 1},
		{Pattern: `^[+-]\s*struct\s*{`, Description: "struct changes", Priority: 2},
		{Pattern: `^[+-]\s*interface\s*{`, Description: "interface changes", Priority: 2},
	},
	"javascript": {
		{Pattern: `^[+-]\s*function\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*const\s+\w+\s*=\s*\([^)]*\)\s*=>`, Description: "arrow functions", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*import\s+`, Description: "import changes", Priority: 2},
	},
	"python": {
		{Pattern: `^[+-]\s*def\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*@\w+`, Description: "decorator changes", Priority: 2},
	},
	"typescript": {
		{Pattern: `^[+-]\s*function\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*export\s+`, Description: "export changes", Priority: 2},
	},
	"java": {
		{Pattern: `^[+-]\s*public\s+class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*public\s+interface\s+\w+`, Description: "interface changes", Priority: 1},
		{Pattern: `^[+-]\s*public\s+enum\s+\w+`, Description: "enum changes", Priority: 1},
		{Pattern: `^[+-]\s*public\s+static\s+void\s+main\s*\(`, Description: "main method changes", Priority: 1},
		{Pattern: `^[+-]\s*public\s+static\s+void\s+\w+\s*\(`, Description: "method changes", Priority: 2},
	},
	"jsx": {
		{Pattern: `^[+-]\s*function\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
	},
	"tsx": {
		{Pattern: `^[+-]\s*function\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
	},
	"swift": {
		{Pattern: `^[+-]\s*func\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*extension\s+\w+`, Description: "extension changes", Priority: 2},
	},
	"kotlin": {
		{Pattern: `^[+-]\s*fun\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*interface\s+\w+`, Description: "interface changes", Priority: 2},
	},
	"php": {
		{Pattern: `^[+-]\s*function\s+\w+`, Description: "function changes", Priority: 1},
		{Pattern: `^[+-]\s*class\s+\w+`, Description: "class changes", Priority: 1},
		{Pattern: `^[+-]\s*interface\s+\w+`, Description: "interface changes", Priority: 2},
	},
	}

type DiffAnalyzer struct {
	MaxLines       int
	ContextLines   int
	ImportantFiles []string
	IgnorePatterns []string
}



// Supported file extensions for code analysis
var relevantExtensions = []string{".go", ".c", ".cpp", ".py", ".js", ".ts", ".tsx", ".jsx"}

// Regex patterns for code analysis
var (
	importRegex   = regexp.MustCompile(`(?i)\b(import|include|from|require)\b.*`)
	functionRegex = regexp.MustCompile(`(?i)\b(func|def|function)\s+[a-zA-Z_][a-zA-Z0-9_]*\s*\(`)
	variableRegex = regexp.MustCompile(`(?i)\b(var|let|const|[a-zA-Z_][a-zA-Z0-9_]*\s*(:?=|:))`)
)

// CodeEssence represents the essential elements extracted from a code file
type CodeEssence struct {
	FilePath  string   `json:"filePath"`
	Imports   []string `json:"imports,omitempty"`
	Functions []string `json:"functions,omitempty"`
	Variables []string `json:"variables,omitempty"`
}

// ExtractEssenceFromFile analyzes a single file and extracts its code essence
func ExtractEssenceFromFile(filePath string) (*CodeEssence, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	essence := &CodeEssence{
		FilePath:  filePath,
		Imports:   make([]string, 0),
		Functions: make([]string, 0),
		Variables: make([]string, 0),
	}

	var multilineComment bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Handle multiline comments
		if strings.HasPrefix(line, "/*") {
			multilineComment = true
			continue
		}
		if multilineComment {
			if strings.Contains(line, "*/") {
				multilineComment = false
			}
			continue
		}

		// Extract code elements
		switch {
		case importRegex.MatchString(line):
			essence.Imports = append(essence.Imports, line)
		case functionRegex.MatchString(line):
			essence.Functions = append(essence.Functions, line)
		case variableRegex.MatchString(line):
			essence.Variables = append(essence.Variables, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file %s: %w", filePath, err)
	}

	return essence, nil
}

// SearchDirectory recursively searches a directory for relevant code files
func SearchDirectory(root string) ([]*CodeEssence, error) {
	var results []*CodeEssence

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}

		if !info.IsDir() && isRelevantFile(path) {
			essence, err := ExtractEssenceFromFile(path)
			if err != nil {
				// Log error but continue processing other files
				fmt.Fprintf(os.Stderr, "Error processing %s: %v\n", path, err)
				return nil
			}
			results = append(results, essence)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", root, err)
	}

	return results, nil
}

// isRelevantFile checks if a file has a supported extension
func isRelevantFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, relevantExt := range relevantExtensions {
		if ext == relevantExt {
			return true
		}
	}
	return false
}

func NewDiffAnalyzer() *DiffAnalyzer {
	return &DiffAnalyzer{
		MaxLines:     2000,
		ContextLines: 3,
		ImportantFiles: []string{
			"main", "core", "api", "service",
			"controller", "model", "repository",
		},
		IgnorePatterns: []string{
			`^\s*//`, `^\s*#`, `^\s*/\*`,
			`^\s*\*`, `^\s*\*/`,
			`^\s*$`,
		},
	}
}

func (da *DiffAnalyzer) analyzeGitDiff(cmd string) (string, error) {
	if cmd == "" {
		return "", fmt.Errorf("empty git diff")
	}

	// splitting into lines for large diffs
	lines := strings.Split(cmd, "\n")
	if len(lines) > da.MaxLines {
		lines = da.filterImportantChanges(lines)
	}

	changes := da.extractMeaningfulChanges(lines)
	summary := da.generateSummary(changes)

	return da.formatCommitMessage(summary), nil
}

func (da *DiffAnalyzer) filterImportantChanges(lines []string) []string {
	var filtered []string
	var language string
	var isImportantFile bool

	for i, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			language = detectLanguage(line)
			isImportantFile = da.isImportantFile(line)
		}

		if isImportantFile && da.isSignificantChange(line, language) {
			// context...
			start := max(0, i-da.ContextLines)
			filtered = append(filtered, lines[start:i+1]...)
		}
	}

	return filtered
}

func (da *DiffAnalyzer) isSignificantChange(line, language string) bool {
	patterns, exists := languagePatterns[language]
	if !exists {
		return false
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern.Pattern, line); matched {
			return true
		}
	}

	return false
}

func (da *DiffAnalyzer) isImportantFile(diffLine string) bool {
	for _, important := range da.ImportantFiles {
		if strings.Contains(diffLine, important) {
			return true
		}
	}
	return false
}

func (da *DiffAnalyzer) extractMeaningfulChanges(lines []string) map[string][]string {
	changes := make(map[string][]string)
	currentFile := ""

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			currentFile = extractFileName(line)
			continue
		}

		if da.isIgnoredLine(line) {
			continue
		}

		if currentFile != "" && (strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")) {
			changes[currentFile] = append(changes[currentFile], line)
		}
	}

	return changes
}

func (da *DiffAnalyzer) generateSummary(changes map[string][]string) string {
	var summaryParts []string

	for file, fileChanges := range changes {
		adds := 0
		removes := 0
		for _, change := range fileChanges {
			if strings.HasPrefix(change, "+") {
				adds++
			} else if strings.HasPrefix(change, "-") {
				removes++
			}
		}

		if adds > 0 || removes > 0 {
			summaryParts = append(summaryParts,
				fmt.Sprintf("%s: +%d/-%d", file, adds, removes))
		}
	}

	return strings.Join(summaryParts, ", ")
}

func (da *DiffAnalyzer) formatCommitMessage(summary string) string {
	template := `
				You are an expert software engineer assisting with writing clear and concise git commit messages. 
				Given the following changes, provide a descriptive commit message in 50 words or less:

				Changes:
				%s

				NOTE: Return only the commit message.
				`
	return fmt.Sprintf(template, summary)
}

func detectLanguage(diffLine string) string {
	if strings.HasSuffix(diffLine, ".go") {
		return "go"
	} else if strings.HasSuffix(diffLine, ".js") || strings.HasSuffix(diffLine, ".ts") {
		return "javascript"
	} else if strings.HasSuffix(diffLine, ".py") {
		return "python"
	} else if strings.HasSuffix(diffLine, ".java") {
		return "java"
	} else if strings.HasSuffix(diffLine, ".c") || strings.HasSuffix(diffLine, ".cpp") {
		return "c"
	} else if strings.HasSuffix(diffLine, ".rb") {
		return "ruby"
	} else if strings.HasSuffix(diffLine, ".cs") {
		return "csharp"
	}
	return ""
}

func extractFileName(diffLine string) string {
	parts := strings.Split(diffLine, " ")
	if len(parts) >= 3 {
		return strings.TrimPrefix(parts[2], "a/")
	}
	return ""
}

func (da *DiffAnalyzer) isIgnoredLine(line string) bool {
	for _, pattern := range da.IgnorePatterns {
		if matched, _ := regexp.MatchString(pattern, line); matched {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
