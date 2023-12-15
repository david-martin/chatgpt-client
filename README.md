# ChatGPT Go Client

## Overview
This application is a simple terminal-based client for interacting with OpenAI's ChatGPT API using Go. It enables users to send messages to the ChatGPT model and receive responses through a command-line interface.

## Prerequisites
- Go 1.20 or later
- An OpenAI API key

## Installation
1. Clone the repository:
   `git clone git@github.com:david-martin/chatgpt-client`

2. Navigate to the project directory:
   `cd chatgpt-client`

3. Build the application (optional):
   `go build`

## Usage
1. Set your OpenAI API key as an environment variable:
   - On Linux/macOS: `export OPENAI_API_KEY=your_api_key_here`
   - On Windows: `set OPENAI_API_KEY=your_api_key_here`

2. Run the application:
   `go run main.go`
   Or, if you built the application:
   `./chatgpt-client`

## Example Output

Below is an example of the interaction with the ChatGPT Go Client:

```text
$ ./chatgpt-client
Enter your message: Hi, can you tell me a joke?

Response: Sure! Why don't scientists trust atoms? Because they make up everything!

Enter your message: That's a good one! Thanks!

Response: You're welcome! If you have any more questions or need another joke, just ask!

Enter your message: exit
```