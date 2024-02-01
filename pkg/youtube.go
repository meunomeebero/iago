package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Comment struct {
	AuthorDisplayName string
	Comments          []string
}

func getComments() []youtube.CommentThreadListResponse {
	ctx := context.Background()

	credentialsFile := "./google.json"

	credentials, err := os.ReadFile(credentialsFile)

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(credentials, youtube.YoutubeForceSslScope)

	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(ctx, config)

	svc, err := youtube.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Unable to create YouTube service: %v", err)
	}

	playlistId := getUploadsPlaylistId(svc)

	videoList, err := getVideosFromPlaylist(svc, playlistId, 30)

	if err != nil {
		log.Fatalf("Error retrieving videos from playlist: %v", err)
	}

	allComments := make([]youtube.CommentThreadListResponse, len(videoList.Items))

	for _, r := range videoList.Items {
		id := ""

		prefix := strings.HasPrefix(r.Id, "VVV")

		if prefix && r.Snippet != nil {

			if r.Snippet.ResourceId != nil {
				id = r.Snippet.ResourceId.VideoId
			}
		}

		if id == "" {
			id = r.Id
		}

		comment, err := getVideoComments(svc, id)

		if err != nil {
			continue
		}

		allComments = append(allComments, *comment)
	}

	return allComments
}

func GetSubscribersComments() []Comment {
	comments := getComments()

	authorToComments := make(map[string][]string)
	formattedComment := make([]Comment, 0)

	for _, c := range comments {
		for _, i := range c.Items {
			authorUsername := i.Snippet.TopLevelComment.Snippet.AuthorDisplayName
			authorToComments[authorUsername] = append(authorToComments[authorUsername], i.Snippet.TopLevelComment.Snippet.TextDisplay)
		}
	}

	for k, v := range authorToComments {
		formattedComment = append(formattedComment, Comment{
			AuthorDisplayName: k,
			Comments:          v,
		})
	}

	return formattedComment
}

func GetChannelComments() []string {
	comments := getComments()

	formattedComments := make([]string, 0)

	for _, c := range comments {
		for _, i := range c.Items {
			formattedComments = append(formattedComments, i.Snippet.TopLevelComment.Snippet.TextDisplay)
		}
	}

	return formattedComments
}

func GetMostSaidWord() (string, int) {
	comments := GetChannelComments()

	words := make(map[string]int)

	wordsToFilter := map[string]bool{
		"que": false,
		" ":   false,
	}

	for _, c := range comments {
		commentWords := strings.Split(c, " ")

		for _, w := range commentWords {
			if len(w) > 1 && wordsToFilter[w] {
				lowerWord := strings.ToLower(w)

				words[lowerWord]++
			}
		}
	}

	var mostSaidWord string
	var timesSaid int

	for k, v := range words {
		if v > timesSaid {
			timesSaid = v
			mostSaidWord = k
		}
	}

	return mostSaidWord, timesSaid
}

func getVideoComments(
	service *youtube.Service,
	videoId string,
) (*youtube.CommentThreadListResponse, error) {
	commentListResponse, err := service.CommentThreads.List([]string{"snippet"}).VideoId(videoId).MaxResults(100).Do()

	if err != nil {
		return nil, err
	}

	return commentListResponse, nil
}

func getUploadsPlaylistId(service *youtube.Service) string {
	channelListResponse, err := service.Channels.List([]string{"contentDetails"}).Mine(true).Do()

	if err != nil {
		log.Fatalf("Error retrieving channel details: %v", err)
	}

	// first playlist is the general playlist
	return channelListResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads
}

func getVideosFromPlaylist(service *youtube.Service, playlistId string, maxResults int64) (*youtube.PlaylistItemListResponse, error) {
	playlistItemsListCall := service.PlaylistItems.List([]string{"snippet"}).PlaylistId(playlistId).MaxResults(maxResults)
	playlistItemsListResponse, err := playlistItemsListCall.Do()

	if err != nil {
		return nil, err
	}

	return playlistItemsListResponse, nil
}

func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	tokFile := "./token.json"

	tok, err := tokenFromFile(tokFile)

	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(ctx, tok)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	tok := &oauth2.Token{}

	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser, then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
