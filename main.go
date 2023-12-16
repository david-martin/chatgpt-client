package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type InternalRequest struct {
	HumanInput     string `json:"human_input"`
	SystemResponse string `json:"system_response"`
}

type InternalResponse struct {
	ResponseToHuman    string `json:"response_to_human"`
	NextKubectlCommand string `json:"next_kubectl_command"`
}

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
	var lastClientResponse string
	lastClientResponse = ""
	for {
		fmt.Print("Enter your message: ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		// Check if the user wants to exit
		if userInput == "exit" {
			fmt.Println("Exiting the application.")
			os.Exit(0) // Exit the program
		}

		// Include last client response in request if available
		internalReq := InternalRequest{
			HumanInput:     userInput,
			SystemResponse: lastClientResponse,
		}
		internalReqJSON, _ := json.Marshal(internalReq)

		var messages []Message
		messages = append(messages, Message{Role: "user", Content: string(internalReqJSON)})

		response, err := getChatGPTResponse(messages)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("Response:", response)
		lastClientResponse = response
	}
}

func getChatGPTResponse(messages []Message) (string, error) {
	// Include your initial instructions along with the messages
	// initialPrompt := `I am using an automated system with a structured message format.
	// Each message has two parts: a human input section and a system response section.
	// The human input section contains commands or questions from a user.
	// The system response section contains outputs from executing Kubernetes commands.
	// Please respond with two sections:
	// 	first, a response to the human input, which can be empty if not necessary;
	// 	second, a Kubernetes command starting with 'kubectl' for the system to execute, which can also be empty if not necessary.
	// Please adhere to this structure in your responses.`
	initialPrompt := `I am using a JSON-based communication protocol for a system that interacts with a Kubernetes cluster. Each message to you will be a JSON object with two fields. The 'human_input' field contains commands or questions from a user about Kubernetes operations. The 'system_response' field contains the output from previously executed Kubernetes commands. Please respond with a JSON object that has two fields: 'response_to_human' and 'next_kubectl_command'. The 'response_to_human' field should contain your direct response to the 'human_input', and the 'next_kubectl_command' field should contain the exact Kubernetes command to be executed next, if any. Ensure that your responses are concise and formatted strictly according to this structure.`

	messages = append([]Message{{Role: "system", Content: initialPrompt}}, messages...)

	url := "https://api.openai.com/v1/chat/completions"
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("API key not set in environment variables")
	}

	// Updated request body
	reqBody, err := json.Marshal(ChatGPTRequest{
		Model:       "gpt-3.5-turbo", // or other model as per your requirement
		Messages:    messages,
		Temperature: 0.7,
	})

	if err != nil {
		return "", err
	}

	// Create a new request
	fmt.Printf("RAW REQUEST: %s\n\n", reqBody)
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
	fmt.Printf("RAW RESPONSE: %v\n\n", chatResponse)

	if len(chatResponse.Choices) > 0 {
		var internalResp InternalResponse
		err := json.Unmarshal([]byte(chatResponse.Choices[0].Message.Content), &internalResp)
		if err != nil {
			fmt.Println("Error parsing response:", err)
		}

		fmt.Println("Response to Human:", internalResp.ResponseToHuman)
		if internalResp.NextKubectlCommand != "" {
			fmt.Println("Executing kubectl command:", internalResp.NextKubectlCommand)
			// Execute kubectl command and update previousSystemResponse
			return executeKubectlCommand(internalResp.NextKubectlCommand), nil
		} else {
			return "", nil
		}
	}

	return "", fmt.Errorf("no response from API")
}

func executeKubectlCommand(command string) string {
	fmt.Println("Executing command:", command)
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Error executing command:", err)
		return "Error: " + err.Error()
	}

	fmt.Println("Command output:", string(output))
	return string(output)
}
