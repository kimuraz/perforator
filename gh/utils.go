package gh

import (
	"context"
	"fmt"
	"time"

	"github.com/aogz/perforator/utils"
	"github.com/google/go-github/v40/github"
)

const MAX_PER_PAGE = 100

// GetPRs returns a list of pull requests filtered by label for specified repo
func GetPRs(client *github.Client, owner string, repo string, limit int, skip int) ([]*github.PullRequest, error) {
	total := limit + skip

	options := &github.PullRequestListOptions{
		State: "closed",
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: limit,
		},
	}

	pages := 1
	var pullRequestList []*github.PullRequest
	if total > MAX_PER_PAGE {
		pages = total / MAX_PER_PAGE
		if total%MAX_PER_PAGE > 0 {
			pages += 1
		}
	}

	for page := 1; page <= pages; page++ {
		utils.ClearPrint(fmt.Sprintf("Retrieving pull requests.. %d/%d", page, pages))
		options.ListOptions.Page = page
		options.ListOptions.PerPage = MAX_PER_PAGE

		prs, _, err := client.PullRequests.List(context.Background(), owner, repo, options)
		if err != nil {
			return pullRequestList, err
		}

		// handle skip
		startSlice := 0
		endSlice := options.ListOptions.PerPage
		if skip > 0 {
			if skip > endSlice {
				skip -= endSlice
				continue
			} else {
				startSlice = skip
				skip = 0
			}
		}

		// calculate remainder in the last page
		if page == pages {
			remainder := total % MAX_PER_PAGE
			if remainder > 0 {
				endSlice = remainder
			}
		}

		pullRequestList = append(pullRequestList, prs[startSlice:endSlice]...)
	}

	return pullRequestList, nil
}

// GetIssuesByRepo returns a list of tickets filtered by label for specified repo
func GetIssuesByRepo(client *github.Client, args utils.DefaultArgs, state string) ([]*github.Issue, error) {
	total := args.Limit + args.Skip

	options := &github.IssueListByRepoOptions{
		State: state,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: args.Limit,
		},
	}

	pages := 1
	var issuesList []*github.Issue
	if total > MAX_PER_PAGE {
		pages = total / MAX_PER_PAGE
		if total%MAX_PER_PAGE > 0 {
			pages += 1
		}
	}

	for page := 1; page <= pages; page++ {
		utils.ClearPrint(fmt.Sprintf("Retrieving issues.. %d/%d", page, pages))
		options.ListOptions.Page = page
		options.ListOptions.PerPage = MAX_PER_PAGE

		issues, _, err := client.Issues.ListByRepo(context.Background(), args.Owner, args.RepoName, options)
		if err != nil {
			return issuesList, err
		}

		// handle skip
		startSlice := 0
		endSlice := options.ListOptions.PerPage
		skip := args.Skip
		if skip > 0 {
			if skip > endSlice {
				skip -= endSlice
				continue
			} else {
				startSlice = skip
				skip = 0
			}
		}

		// calculate remainder in the last page
		if page == pages {
			remainder := total % MAX_PER_PAGE
			if remainder > 0 {
				endSlice = remainder
			}
		}

		issuesList = append(issuesList, issues[startSlice:endSlice]...)
	}

	return issuesList, nil
}

// GetPullRequestTimeline ...
func GetPullRequestTimeline(client *github.Client, pr *github.PullRequest) ([]*github.Timeline, error) {
	options := &github.ListOptions{
		PerPage: MAX_PER_PAGE,
	}
	timeline, _, err := client.Issues.ListIssueTimeline(context.Background(), *pr.Head.Repo.Owner.Login, *pr.Head.Repo.Name, *pr.Number, options)
	return timeline, err
}

func GetPullRequestReviews(client *github.Client, owner string, repo string, prNumber int) ([]*github.PullRequestReview, error) {
	options := &github.ListOptions{
		PerPage: MAX_PER_PAGE,
	}
	reviews, _, err := client.PullRequests.ListReviews(context.Background(), owner, repo, prNumber, options)
	return reviews, err
}

func GetPullRequest(client *github.Client, owner string, repo string, prNumber int) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Get(context.Background(), owner, repo, prNumber)
	return pr, err
}

func GetCommits(client *github.Client, owner string, repo string, author string, since, until time.Time) ([]*github.RepositoryCommit, error) {
	options := &github.CommitsListOptions{
		Author: author,
		Since:  since,
		Until:  until,
	}
	commits, _, err := client.Repositories.ListCommits(context.Background(), owner, repo, options)
	return commits, err
}
