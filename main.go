package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message OpenAIMessage `json:"message"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("english> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("Goodbye!")
			break
		}

		out := evaluate(input)
		fmt.Println(out)
	}
}

func evaluate(input string) string {
	verboseMode := false
	answerOnlyMode := false

	if strings.HasPrefix(input, "!") {
		verboseMode = true
	}

	if strings.HasPrefix(input, "?") {
		answerOnlyMode = true
	}

	exCmd, err := getCommand(input)
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	if answerOnlyMode {
		return exCmd
	}

	if verboseMode {
		fmt.Println(":>", exCmd)
		// wait for user to press enter then continue
		fmt.Print("Press enter to continue...")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}

	cmd := exec.Command("sh", "-c", exCmd)

	out, err := cmd.Output()
	if err != nil {
		log.Println("err executing command: ", err)
		return fmt.Sprintf("Error: %s", err)
	}
	return string(out)
}

func getCommand(data string) (string, error) {

	// Set your OpenAI API key from environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")

	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable not set")
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	requestData := OpenAIRequest{
		Model: "gpt-4o", // Example model, replace with the one you want to use
		Messages: []OpenAIMessage{
			{Role: "system", Content: "Convert given sentence into an executable shell command, return only the command."},
			{Role: "user", Content: data},
		},
		Temperature: 0.7,
		MaxTokens:   100,
	}

	// Marshal the request data to JSON
	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling request data:", err)
		return "", err
	}

	// Create the HTTP POST request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	// Parse the response JSON
	var openAIResponse OpenAIResponse
	err = json.Unmarshal(body, &openAIResponse)
	if err != nil {
		fmt.Println("Error parsing response JSON:", err)
		return "", err
	}

	executableCommand := ""
	if len(openAIResponse.Choices) > 0 {
		executableCommand = openAIResponse.Choices[0].Message.Content
	} else {
		fmt.Println("No executable command found in the response.")
	}

	command := strings.TrimPrefix(executableCommand, "```sh\n")
	command = strings.TrimSuffix(command, "\n```")
	// log.Println("cmd is ", command)
	return command, nil

}
