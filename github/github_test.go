package github

import (
	"context"
	"testing"

	"github.com/google/go-github/github"
)

// Testing strategy found here: https://github.com/google/go-github/issues/113#issuecomment-46023864
type MockSearchService struct {
	name string
}

func (s *MockSearchService) Users(ctx context.Context, query string, opt *github.SearchOptions) (*github.UsersSearchResult, *github.Response, error) {
	user := s.name
	users := []github.User{github.User{Login: &user}}
	var resp *github.Response

	results := &github.UsersSearchResult{
		// We need pointers to values, so use an anonymous function.
		Total:             func() *int { i := int(1); return &i }(),
		IncompleteResults: func() *bool { i := bool(false); return &i }(),
		Users:             users,
	}

	return results, resp, nil
}

func TestGetUsernameForEmail(t *testing.T) {
	Client = GitHubClient{Search: &MockSearchService{name: "Jason"}}

	username, _ := GetUsernameForEmail("fake@email.dev")

	if username != "Jason" {
		t.Errorf("Expecting Jason, got %s", username)
	}
}
