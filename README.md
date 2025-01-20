# 📚 README

## 👋 Introduction

This project is a Go-based command-line interface (CLI) that utilizes the Groq API to provide a simple and efficient way to generate descriptive commit messages and chat responses. 🤖

## 🚀 Project Setup

### Prerequisites

1. **Install Go**: Download and install Go from the official website: https://go.dev/dl/ 📦

2. **Install Make**:
   - **Windows**:
     ```bash
     # Using Chocolatey
     choco install make
     # or using Scoop
     scoop install make
     ```
   - **macOS**:
     ```bash
     # Using Homebrew
     brew install make
     ```
   - **Linux**:
     ```bash
     # Debian/Ubuntu
     sudo apt-get install make
     # Fedora
     sudo dnf install make
     # CentOS/RHEL
     sudo yum install make
     ```

3. **Get the Repository**: Clone the repository:
   ```bash
   git clone https://github.com/NevroHelios/cliagent.git
   ```

4. **Install Dependencies**:
   ```bash
   go get -u github.com/joho/godotenv
   go get -u github.com/kelvins/go-uber-zap
   go get -u github.com/urfave/cli/v2
   ```

### Building the Project

#### Development Build (using .env)
1. Create a `.env` file in the project root:
   ```makefile
   GROQ_API_KEY=your_api_key_here
   ```
2. Build the project:
   ```bash
   go build -o gocli
   ```

#### Production Build (embedded API key)
Build the project with your API key embedded:

```bash
# Linux/macOS
go build -ldflags "-X main.GROQ_API_KEY=your-api-key-here" -o gocli

# Windows (PowerShell)
go build -ldflags "-X main.GROQ_API_KEY=your-api-key-here" -o gocli.exe
```

Alternatively, use environment variable:
```bash
# Linux/macOS
export GROQ_API_KEY=your-api-key-here
go build -ldflags "-X main.GROQ_API_KEY=$GROQ_API_KEY" -o gocli

# Windows (PowerShell)
$env:GROQ_API_KEY="your-api-key-here"
go build -ldflags "-X main.GROQ_API_KEY=$env:GROQ_API_KEY" -o gocli.exe
```

### Creating a Symlink (Optional)

To run the CLI from any directory:

- **Windows** (Run PowerShell as Administrator):
  ```powershell
  New-Item -ItemType SymbolicLink -Path "C:\Windows\System32\gocli.exe" -Target "$pwd\gocli.exe"
  ```

- **Linux/macOS**:
  ```bash
  sudo ln -s "$(pwd)/gocli" /usr/local/bin/gocli
  ```

## 🏃‍♂️ Running the Project

1. **Navigate to the Project Directory**:
   ```bash
   cd cliagent
   ```

2. **Run the Project**:
   ```bash
   ./gocli
   # If symlinked, simply run from anywhere:
   gocli
   ```

3. **Follow the Prompts**: Follow the prompts to select a model, enter a query, and generate a response. 🤔

## 📝 Generating Commit Messages

1. **Select the "Commit" Option**: Select the "Commit" option from the main menu.
2. **Select a Model**: Enter a query to generate a commit message.
3. **Generate the Commit Message**: The project will generate a descriptive commit message with your selected Model. 📝

## 🤖 Generating Chat Responses

1. **Select the "Chat" Option**: Select the "Chat" option from the main menu.
2. **Enter a Query**: Enter a query to generate a chat response.
3. **Generate the Chat Response**: The project will generate a chat response based on your query. 💬

## 👍 Contributing

If you'd like to contribute to the project, please follow these steps:

1. **Fork the Repository**: Fork the repository using the GitHub fork button.
2. **Make Changes**: Make changes to the code and commit them using a descriptive commit message.
3. **Create a Pull Request**: Create a pull request to merge your changes into the main repository. 🚀

## ⚠️ Security Note

When building with an embedded API key:
- Never commit the build command containing your API key
- Consider using environment variables for the build process
- Keep your compiled binary secure as it contains your API key

## 🙏 Acknowledgments

This project utilizes the following libraries and frameworks:

* Go: https://go.dev/
* Groq API: https://api.groq.com/
* Godotenv: https://github.com/joho/godotenv
* Uber-zap: https://github.com/urfave/cli/v2

Thanks to the maintainers and contributors of these libraries and frameworks for their hard work and dedication! 🙏