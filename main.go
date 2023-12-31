package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGPTRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type ChatGPTResponse struct {
	Choices []Choice `json:"choices"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter your message: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		// Check if the user wants to exit
		if userInput == "exit" {
			fmt.Println("Exiting the application.")
			os.Exit(0) // Exit the program
		}
		response, err := getChatGPTResponse(userInput)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		fmt.Println("Response:", response)

		// To break the loop, you might want to check for a specific command like "exit".
	}
}

func getChatGPTResponse(input string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set in environment variables")
	}

	// Updated request body
	reqBody, err := json.Marshal(ChatGPTRequest{
		Model: "gpt-3.5-turbo", // or other model as per your requirement
		Messages: []Message{
			{
				Role:    "user",
				Content: input,
			},
		},
		Temperature: 0.7,
	})

	if err != nil {
		return "", err
	}

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	// Set the required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Read the response
	var chatResponse ChatGPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResponse); err != nil {
		return "", err
	}

	if len(chatResponse.Choices) > 0 {
		return chatResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from API")
}
