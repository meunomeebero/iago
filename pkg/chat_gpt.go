package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil" // Import this package to read the response body.
	"net/http"

	"github.com/robertokbr/iago/utils"
)

type _ChatAnswer struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type RequestPayloadMessages struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type RequestPayload struct {
	Model    string                   `json:"model"`
	Messages []RequestPayloadMessages `json:"messages"`
}

func AnswerQuestion(prompts ...string) string {
	// Create the payload using a struct.
	payload := RequestPayload{
		Model: "gpt-4",
		Messages: []RequestPayloadMessages{
			{
				Role: "system",
			},
		},
	}

	for _, r := range prompts {
		payload.Messages = append(payload.Messages, RequestPayloadMessages{
			Role:    "user",
			Content: r,
		})
	}

	fp, _ := json.Marshal(payload)

	utils.CreateJSONFile("comments.json", fp)

	// Marshal the payload into JSON.
	data, err := json.Marshal(payload)

	utils.CreateJSONFile("file.json", data)

	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return ""
	}

	baseURL := "https://api.openai.com/v1/chat/completions"
	request, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(data))
	if err != nil {
		fmt.Printf("Error creating request: %v", err)
		return ""
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+utils.GetSK())

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body) // Read the response body for error information.
		fmt.Printf("Error response from API: Status Code %d, Response: %s\n", resp.StatusCode, string(bodyBytes))
		return ""
	}

	var res _ChatAnswer

	// Check if there's an error when decoding the response.
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		fmt.Printf("Error decoding API response: %v", err)
		return ""
	}

	// Assuming there's at least one choice and it has content.
	if len(res.Choices) > 0 {
		return res.Choices[0].Message.Content
	}

	return ""
}
