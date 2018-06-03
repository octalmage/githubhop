package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
)

type searchService interface {
	Users(ctx context.Context, query string, opt *github.SearchOptions) (*github.UsersSearchResult, *github.Response, error)
}

// GitHubClient Type for go-github client.
type GitHubClient struct {
	Search searchService
}

func newClient(httpClient *http.Client) GitHubClient {
	client := github.NewClient(httpClient)
	// optionally set client.BaseURL, client.UserAgent, etc

	return GitHubClient{
		Search: client.Search,
	}
}

// Client The GitHub client used by GetUsernameForEmail.
var Client = newClient(nil)

// GetUsernameForEmail Lookup GitHub username using an email address.
func GetUsernameForEmail(email string) (string, error) {
	opts := &github.SearchOptions{Order: "desc", ListOptions: github.ListOptions{Page: 1, PerPage: 1}}
	query := fmt.Sprintf("%s in:email", email)
	result, _, err := Client.Search.Users(context.Background(), query, opts)

	if err != nil {
		return "", err
	}

	return *result.Users[0].Login, nil
}
