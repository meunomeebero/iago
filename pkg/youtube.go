package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	// ClientID     string       `json:"client_id"`
	// ClientSecret string       `json:"client_secret"`
	// AccessToken  string       `json:"access_token"`
	// RefreshToken string       `json:"refresh_token"`
	// Scope        string       `json:"scope"`
	// TokenType    string       `json:"token_type"`
	// Expiry       oauth2.Token `json:"expiry_date"`
}

func (self *YouTube) LoadData() {
	// config := oauth2.Config{
	// 	ClientID: self.ClientID,
	// 	ClientSecret: self.ClientSecret,
	// 	Endpoint: google.Endpoint,
	// 	RedirectURL: "urn:ietf:wg:oauth:2.0:oob",
	// 	Scopes: []string{youtube.YoutubeReadonlyScope},
	// }

	// token := oauth2.Token{
	// 	AccessToken: self.AccessToken,
	// 	RefreshToken: self.RefreshToken,
	// 	TokenType: self.TokenType,
	// 	Expiry: self.Expiry.Expiry,
	// }

	ctx := context.Background()

	svc, err := youtube.NewService(ctx, option.WithCredentialsFile("./google-token.json"))

	if err != nil {
		return
	}

	res, err := svc.Videos.List([]string{"snippet", "contentDetails"}).Do()

	if err != nil {
		return
	}

	file, _ := os.Create("res.json")

	defer file.Close()

	data, err := json.Marshal(res.Items)

	jsonData := bytes.NewBuffer(data)

	io.Copy(file, jsonData)
}

// func NewYouTube() (*YouTube, error) {
// 	data, err := os.ReadFile("./google-token.json")

// 	if err != nil {
// 		return &YouTube{}, err
// 	}

// 	var ytb YouTube

// 	err = json.Unmarshal(data, &ytb)

// 	return &ytb, nil
// }
