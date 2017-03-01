package gitstats

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubClient struct {
	Client *github.Client
}

func NewClient(accessToken string) *GithubClient {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	ghc := GithubClient{}
	ghc.Client = github.NewClient(tc)

	return &ghc
}

func (ghc *GithubClient) GetUsers(ctx context.Context, user string) (*github.User, *github.Response, error) {
	return ghc.Client.Users.Get(ctx, user)
}

func (ghc *GithubClient) GetOrganizations(ctx context.Context, org string) (*github.Organization, *github.Response, error) {
	return ghc.Client.Organizations.Get(ctx, org)
}

func (ghc *GithubClient) GetRepository(ctx context.Context, owner string, repo string) (*github.Repository, *github.Response, error) {
	return ghc.Client.Repositories.Get(ctx, owner, repo)
}

func (ghc *GithubClient) ListRepositories(ctx context.Context, owner string, opt *github.RepositoryListOptions) ([]*github.Repository, *github.Response, error) {
	return ghc.Client.Repositories.List(ctx, owner, opt)
}
