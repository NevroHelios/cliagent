📚 README
================

👋 Introduction
---------------

This project is a Go-based command-line interface (CLI) that utilizes the Groq API to provide a simple and efficient way to generate descriptive commit messages and chat responses. 🤖

🚀 Project Setup
-----------------

To set up the project, follow these steps:

1. **Install Go**: Make sure you have Go installed on your system. You can download it from the official Go website: https://go.dev/dl/ 📦
2. **Get the Repository**: Clone the repository using the following command:
```bash
git clone https://github.com/NevroHelios/cliagent.git
```
3. **Install Dependencies**: Install the required dependencies using the following command:
```bash
go get -u github.com/joho/godotenv
go get -u github.com/kelvins/go-uber-zap
go get -u github.com/urfave/cli/v2
```
4. **Create a `.env` File**: Create a new file named `.env` in the root of the project and add your Groq API key:
```makefile
GROQ_API_KEY=your_api_key_here
```
🏃‍♂️ Running the Project
-------------------------

To run the project, follow these steps:

1. **Navigate to the Project Directory**: Navigate to the project directory using the following command:
```bash
cd cliagent
```
2. **Build the Project**: Build the project using the following command:
```bash
go build -o gocli
```
3. **Run the Project**: Run the project using the following command:
```bash
./gocli
```
4. **Follow the Prompts**: Follow the prompts to select a model, enter a query, and generate a response. 🤔

📝 Generating Commit Messages
-----------------------------

To generate a descriptive commit message, follow these steps:

1. **Select the "Commit" Option**: Select the "Commit" option from the main menu.
2. **Select a Model**: Enter a query to generate a commit message.
3. **Generate the Commit Message**: The project will generate a descriptive commit message with your selected Model. 📝

🤖 Generating Chat Responses
-----------------------------

To generate a chat response, follow these steps:

1. **Select the "Chat" Option**: Select the "Chat" option from the main menu.
2. **Enter a Query**: Enter a query to generate a chat response.
3. **Generate the Chat Response**: The project will generate a chat response based on your query. 💬

👍 Contributing
-----------------

If you'd like to contribute to the project, please follow these steps:

1. **Fork the Repository**: Fork the repository using the GitHub fork button.
2. **Make Changes**: Make changes to the code and commit them using a descriptive commit message.
3. **Create a Pull Request**: Create a pull request to merge your changes into the main repository. 🚀

🙏 Acknowledgments
------------------

This project utilizes the following libraries and frameworks:

* Go: https://go.dev/
* Groq API: https://api.groq.com/
* Godotenv: https://github.com/joho/godotenv
* Uber-zap: https://github.com/urfave/cli/v2

Thanks to the maintainers and contributors of these libraries and frameworks for their hard work and dedication! 🙏