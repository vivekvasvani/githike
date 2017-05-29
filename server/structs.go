package server

import "time"

type SearchReposStruct struct {
	TotalCount        int  `json:"total_count"`
	IncompleteResults bool `json:"incomplete_results"`
	Items             []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Owner    struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"owner"`
		Private          bool        `json:"private"`
		HTMLURL          string      `json:"html_url"`
		Description      string      `json:"description"`
		Fork             bool        `json:"fork"`
		URL              string      `json:"url"`
		ForksURL         string      `json:"forks_url"`
		KeysURL          string      `json:"keys_url"`
		CollaboratorsURL string      `json:"collaborators_url"`
		TeamsURL         string      `json:"teams_url"`
		HooksURL         string      `json:"hooks_url"`
		IssueEventsURL   string      `json:"issue_events_url"`
		EventsURL        string      `json:"events_url"`
		AssigneesURL     string      `json:"assignees_url"`
		BranchesURL      string      `json:"branches_url"`
		TagsURL          string      `json:"tags_url"`
		BlobsURL         string      `json:"blobs_url"`
		GitTagsURL       string      `json:"git_tags_url"`
		GitRefsURL       string      `json:"git_refs_url"`
		TreesURL         string      `json:"trees_url"`
		StatusesURL      string      `json:"statuses_url"`
		LanguagesURL     string      `json:"languages_url"`
		StargazersURL    string      `json:"stargazers_url"`
		ContributorsURL  string      `json:"contributors_url"`
		SubscribersURL   string      `json:"subscribers_url"`
		SubscriptionURL  string      `json:"subscription_url"`
		CommitsURL       string      `json:"commits_url"`
		GitCommitsURL    string      `json:"git_commits_url"`
		CommentsURL      string      `json:"comments_url"`
		IssueCommentURL  string      `json:"issue_comment_url"`
		ContentsURL      string      `json:"contents_url"`
		CompareURL       string      `json:"compare_url"`
		MergesURL        string      `json:"merges_url"`
		ArchiveURL       string      `json:"archive_url"`
		DownloadsURL     string      `json:"downloads_url"`
		IssuesURL        string      `json:"issues_url"`
		PullsURL         string      `json:"pulls_url"`
		MilestonesURL    string      `json:"milestones_url"`
		NotificationsURL string      `json:"notifications_url"`
		LabelsURL        string      `json:"labels_url"`
		ReleasesURL      string      `json:"releases_url"`
		DeploymentsURL   string      `json:"deployments_url"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
		PushedAt         time.Time   `json:"pushed_at"`
		GitURL           string      `json:"git_url"`
		SSHURL           string      `json:"ssh_url"`
		CloneURL         string      `json:"clone_url"`
		SvnURL           string      `json:"svn_url"`
		Homepage         interface{} `json:"homepage"`
		Size             int         `json:"size"`
		StargazersCount  int         `json:"stargazers_count"`
		WatchersCount    int         `json:"watchers_count"`
		Language         string      `json:"language"`
		HasIssues        bool        `json:"has_issues"`
		HasProjects      bool        `json:"has_projects"`
		HasDownloads     bool        `json:"has_downloads"`
		HasWiki          bool        `json:"has_wiki"`
		HasPages         bool        `json:"has_pages"`
		ForksCount       int         `json:"forks_count"`
		MirrorURL        interface{} `json:"mirror_url"`
		OpenIssuesCount  int         `json:"open_issues_count"`
		Forks            int         `json:"forks"`
		OpenIssues       int         `json:"open_issues"`
		Watchers         int         `json:"watchers"`
		DefaultBranch    string      `json:"default_branch"`
		Score            int         `json:"score"`
	} `json:"items"`
}

type SlackAppRequest struct {
	Actions []struct {
		Name            string `json:"name"`
		Type            string `json:"type"`
		SelectedOptions []struct {
			Value string `json:"value"`
		} `json:"selected_options"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	IsAppUnfurl     bool   `json:"is_app_unfurl"`
	OriginalMessage struct {
		Text        string `json:"text"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID         string `json:"id"`
				Name       string `json:"name"`
				Text       string `json:"text"`
				Type       string `json:"type"`
				DataSource string `json:"data_source,omitempty"`
				Options    []struct {
					Text  string `json:"text"`
					Value string `json:"value"`
				} `json:"options,omitempty"`
				Value string `json:"value,omitempty"`
				Style string `json:"style,omitempty"`
			} `json:"actions"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
}

type RequestButtonStruct struct {
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"actions"`
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID     string `json:"id"`
		Domain string `json:"domain"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs        string `json:"action_ts"`
	MessageTs       string `json:"message_ts"`
	AttachmentID    string `json:"attachment_id"`
	Token           string `json:"token"`
	IsAppUnfurl     bool   `json:"is_app_unfurl"`
	OriginalMessage struct {
		Text        string `json:"text"`
		BotID       string `json:"bot_id"`
		Attachments []struct {
			CallbackID string `json:"callback_id"`
			Text       string `json:"text"`
			ID         int    `json:"id"`
			Color      string `json:"color"`
			Actions    []struct {
				ID    string `json:"id"`
				Name  string `json:"name"`
				Text  string `json:"text"`
				Type  string `json:"type"`
				Value string `json:"value"`
				Style string `json:"style"`
			} `json:"actions"`
		} `json:"attachments"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
	} `json:"original_message"`
	ResponseURL string `json:"response_url"`
}

type TeamAdminMapStr struct {
	Teamadminmap []struct {
		Teamname string `json:"teamname"`
		Admin    string `json:"admin"`
	} `json:"teamadminmap"`
}

type StatusMessage struct {
	StatusCode string `json:"statusCode"`
	StatusType string `json:"statusType"`
	Message    string `json:"message"`
}

//Response structure
type Response struct {
	Status StatusMessage `json:"status"`
	Data   interface{}   `json:"data,omitempty"`
}
