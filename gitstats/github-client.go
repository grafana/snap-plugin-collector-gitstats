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

func (ghc *GithubClient) GetAllIssues(ctx context.Context, owner string, repo string) ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{State: "all"}
	opt.ListOptions = github.ListOptions{PerPage: 100}
	var allIssues []*github.Issue

	for {
		issues, resp, err := ghc.Client.Issues.ListByRepo(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
	}

	return allIssues, nil
}

func (ghc *GithubClient) GetAllLabels(ctx context.Context, owner string, repo string) ([]*github.Label, error) {
	opt := &github.ListOptions{PerPage: 100}
	var allLabels []*github.Label

	for {
		labels, resp, err := ghc.Client.Issues.ListLabels(ctx, owner, repo, opt)
		if err != nil {
			return nil, err
		}
		allLabels = append(allLabels, labels...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allLabels, nil
}

func (ghc *GithubClient) GetAllLabelsAndIssues(ctx context.Context, owner string, repo string) ([]*github.Label, []*github.Issue, error) {
	labels, err := ghc.GetAllLabels(ctx, owner, repo)
	if err != nil {
		return nil, nil, err
	}

	issues, err := ghc.GetAllIssues(ctx, owner, repo)
	if err != nil {
		return nil, nil, err
	}

	return labels, issues, nil
}
