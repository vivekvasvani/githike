package githubapis

import (
	"context"
	"strconv"
	"time"

	"fmt"

	"encoding/json"

	"github.com/google/go-github/github"
	"github.com/vivekvasvani/githike/client"
	"golang.org/x/oauth2"
)

const (
	GITHUB_TOKEN = ""
	ORG          = "hike"
)

//Check is user is valid
func CheckIfUserIsMemberOfOrg(userid string) (bool, string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	response, _, err := client.Organizations.IsMember(ctx, ORG, userid)
	if err != nil {
		return false, "Error in getting response"
	}
	return response, "nil"
}

//Get Team names and ids map
func GetTeamNamesAndIdsMap() map[string]int {
	nameIdMap := make(map[string]int)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opts := &github.ListOptions{
		PerPage: 200,
	}
	teams, _, _ := client.Organizations.ListTeams(ctx, ORG, opts)
	for _, val := range teams {
		nameIdMap[val.GetName()] = val.GetID()
	}
	return nameIdMap
}

//List Team Members
func ListTeamMembers(teamName, role string) ([]*github.User, string) {
	if role == "" {
		role = "all"
	}
	opts := &github.OrganizationListTeamMembersOptions{
		Role: role,
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	users, _, err := client.Organizations.ListTeamMembers(ctx, GetTeamNamesAndIdsMap()[teamName], opts)
	if err != nil {
		return nil, "Error in fetching Team Members"
	}
	return users, ""
}

//List all pull requests
func ListPullRequests(repo, state, sort, direction string) ([]*github.PullRequest, string) {
	var pullrequests []*github.PullRequest
	if state == "" {
		state = "open"
	}
	if sort == "" {
		sort = "created"
	}
	if direction == "" {
		direction = "desc"
	}
	opts := &github.PullRequestListOptions{
		State:     state,
		Sort:      sort,
		Direction: direction,
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	/*
		repository, _, err := client.Repositories.Get(ctx, ORG, repo)
		if err != nil {
			fmt.Println(err)
			return nil, "Error in fetching repository"
		}
	*/
	//if repository.GetName() == repo {
	pullrequests, _, errinfetch := client.PullRequests.List(ctx, ORG, repo, opts)
	if errinfetch != nil {
		fmt.Println(errinfetch)
		return nil, "Error in fetching pull requests"
	}
	//}
	return pullrequests, ""
}

func ListTeams() ([]*github.Team, string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	opts := &github.ListOptions{
		PerPage: 200,
	}
	teams, _, err := client.Organizations.ListTeams(ctx, ORG, opts)
	if err != nil {
		return nil, "Error in fetching teams info"
	}
	return teams, ""
}

func ListRepos() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, ORG, opt)
		if err != nil {
			fmt.Println(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
}

func AddTeamMembership(team int, user, role string) (*github.Membership, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)

	opts := &github.OrganizationAddTeamMembershipOptions{
		Role: role,
	}
	client := github.NewClient(tc)
	membership, _, err := client.Organizations.AddTeamMembership(ctx, team, user, opts)
	return membership, err
}

func SearchRepos() []string {
	header := make(map[string]string)
	var (
		searchStruct SearchReposStruct
		repos        []string
	)
	header["Authorization"] = "token " + GITHUB_TOKEN
	url := "https://api.github.com/search/repositories?per_page=50&q=user:hike+pushed:>" + time.Now().AddDate(0, -2, 0).Format("2006-01-02")
	response := client.HitRequest(url, "GET", header, "")
	err := json.Unmarshal(response, &searchStruct)
	fmt.Println(err)
	for _, val := range searchStruct.Items {
		repos = append(repos, val.Name)
	}
	return repos
}

func DeactivateGithubHikeAccount(githubid string) bool {
	var result = true
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	_, err := client.Organizations.RemoveMember(ctx, githubid, ORG)
	if err != nil {
		result = false
	}
	return result
}

func GetGithubIdFromEmail(email string) string {
	header := make(map[string]string)
	var (
		searchStruct GetGithubId
	)
	header["Authorization"] = "token " + GITHUB_TOKEN
	url := "https: //api.github.com/search/users?q=" + email + "+in:email+type:users"
	response := client.HitRequest(url, "GET", header, "")
	err := json.Unmarshal(response, &searchStruct)
	fmt.Println(err)
	return searchStruct.Items[0].Login
}

func CreateRepository(name, description, private string, teamid int) bool {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(ctx, ts)
	privateBool, _ := strconv.ParseBool(private)
	opts := &github.Repository{
		Name:        &name,
		Description: &description,
		TeamID:      &teamid,
		Private:     &privateBool,
	}
	client := github.NewClient(tc)
	_, response, err := client.Repositories.Create(ctx, ORG, opts)
	if response.StatusCode == 201 && err == nil {
		return true
	}
	return false
}
