package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
)

// GetUsernameForEmail Lookup GitHub username using an email address.
func GetUsernameForEmail(email string) (string, error) {
	client := github.NewClient(nil)
	opts := &github.SearchOptions{Order: "desc", ListOptions: github.ListOptions{Page: 1, PerPage: 1}}
	query := fmt.Sprintf("%s in:email", email)
	result, _, err := client.Search.Users(context.Background(), query, opts)

	if err != nil {
		return "", err
	}

	return *result.Users[0].Login, nil
}
