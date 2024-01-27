package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/robertokbr/iago/utils"
)

type _ChatAnswer struct {
	Choices []struct {
		Message struct {
			Content string
		}
	}
}

func AnswerQuestion(prompt string) string {
	data := `
		{
			"model": "gpt-3.5-turbo",
			"messages": [
				{
					"role": "system",
					"content": "%s"
				}
			]
		}
	`

	data = fmt.Sprintf(data, prompt)
	dataB := bytes.NewBuffer([]byte(data))
	baseURL := "https://api.openai.com/v1/chat/completions"
	client, err := http.NewRequest("POST", baseURL, dataB)
	client.Header.Add("Content-Type", "application/json")
	client.Header.Add("Authorization", "Bearer "+utils.GetSK())
	resp, err := http.DefaultClient.Do(client)

	if err != nil {
		fmt.Printf("%v", err)
		return ""
	}

	defer resp.Body.Close()

	var res _ChatAnswer

	json.NewDecoder(resp.Body).Decode(&res)

	return res.Choices[0].Message.Content
}
