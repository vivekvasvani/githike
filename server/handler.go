package server

import (
	"database/sql"
	"encoding/json"

	"fmt"

	"strconv"

	"strings"

	"net/http"

	"github.com/nlopes/slack"
	"github.com/valyala/fasthttp"
	"github.com/vivekvasvani/githike/client"
	git "github.com/vivekvasvani/githike/githubapis"
	slackapis "github.com/vivekvasvani/githike/slackapis"
)

const (
	application_json               = "application/json"
	SLACK_TOKEN                    = ""
	SLACK_WEBHOOK                  = ""
	SLACK_WEBHOOK_TO_SEND_SLACKBOT = ""
)

var (
	header  = make(map[string]string)
	output  = make([]string, 0)
	release = make([]string, 0)
	listPRs = make(map[string]string)
)

func SendGitHikeOptions(ctx *fasthttp.RequestCtx) {
	//headers for response
	header["Content-Type"] = application_json
	header["Accept"] = application_json
	//Send 200 OK response immidiately
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(http.StatusOK)

	//Send back the response/options to user
	client.HitRequest(string(ctx.PostArgs().Peek("response_url")), "POST", header, GetPayload("gitoptions.json"))

}

func HandleAppRequests(ctx *fasthttp.RequestCtx) {
	//Send 200 OK response immidiately
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.SetStatusCode(http.StatusOK)

	slackapis.InviteUserToChannel()
	var appRequest SlackAppRequest
	err := json.Unmarshal(ctx.Request.PostArgs().Peek("payload"), &appRequest)
	if err != nil {
		fmt.Println(err)
	}

	header["Content-Type"] = application_json
	header["Accept"] = application_json
	switch appRequest.Actions[0].Type {
	case "select":
		switch appRequest.Actions[0].SelectedOptions[0].Value {
		case "ListPRs":
			var repos = make([]string, 2)
			var session = make([]string, 1)
			//var wg sync.WaitGroup
			var (
				response_url string
				options      string
			)
			response_url = appRequest.ResponseURL
			client.HitRequest(response_url, "POST", header, "{ \"text\": \"Wait... Fetching all repos!!!\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
			func() {
				repos = git.SearchRepos()
			}()
			client.HitRequest(response_url, "POST", header, "{ \"text\": \"Yeah!!! got all repos!!!\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
			for _, val := range repos {
				options = options + "{\"text\": \"" + val + "\",\"value\": \"" + val + "\"},"
			}
			session[0] = options[0 : len(options)-1]
			payload := SubstParams(session, GetPayload("sendListPROptions.json"))
			client.HitRequest(response_url, "POST", header, payload)

		case "AddUserToTeam":
			response_url := appRequest.ResponseURL
			var (
				session                = make([]string, 2)
				valuesForTeamsDropDown string
				roleDefinition         string = "Role specifies the role the user should have in the team. Possible values are:\n" +
					"1. member - a normal member of the team\n" +
					"2. maintainer - a team maintainer. Able to add/remove other team members, promote other team members to team maintainer, and edit the team’s name and description"
			)
			client.HitRequest(response_url, "POST", header, "{ \"text\": \"Wait... Fetching all Teams!!!\", \"response_type\": \"in_channel\", \"replace_original\": true }")
			allTeamsArray, _ := git.ListTeams()
			for i, val := range allTeamsArray {
				valuesForTeamsDropDown = valuesForTeamsDropDown + "{ \"title\": \"\", \"value\": \"" + strconv.Itoa(i+1) + ". " + val.GetName() + "\", \"short\": true },"
			}
			session[0] = valuesForTeamsDropDown[0 : len(valuesForTeamsDropDown)-1]
			session[1] = roleDefinition
			client.HitRequest(response_url, "POST", header, "{ \"text\": \"Wait... Processing your request!!!\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
			payload := SubstParams(session, GetPayload("addusertoteam.json"))
			client.HitRequest(response_url, "POST", header, payload)

		case "DeleteUser":
			//response_url := appRequest.ResponseURL

		}

		switch appRequest.CallbackID {
		case "list_pr_options":
			switch appRequest.Actions[0].Name {
			case "select_repo":
				//repo =
				listPRs["repo"] = appRequest.Actions[0].SelectedOptions[0].Value
			case "select_state":
				var sendPR string
				response_url := appRequest.ResponseURL
				prefix := "#PR\t|\t\t\t\t\t\t\t\t\t\t\tTitle\t\t\t\t\t\t\t\t\t\t\t|\t\t\tcreated_at\t\t\t|\tstatus\n"
				state := appRequest.Actions[0].SelectedOptions[0].Value
				client.HitRequest(response_url, "POST", header, "{ \"text\": \"Your request has been submitted!!!\", \"response_type\": \"in_channel\", \"replace_original\": true }")
				pulls, err := git.ListPullRequests(listPRs["repo"], state, "", "")
				if err != "" {
					fmt.Println(err)
				}
				for _, va := range pulls {
					sendPR = sendPR + "<" + va.GetHTMLURL() + " | " + strconv.Itoa(va.GetNumber()) + ">\t|\t" + va.GetTitle() + "\t|\t" + va.GetCreatedAt().String() + "\t|\t" + strings.ToUpper(va.GetState()) + "\n"
				}
				finalResp := prefix + sendPR
				client.HitRequest(response_url, "POST", header, "{\"text\" : \""+finalResp+"\", \"response_type\": \"in_channel\", \"replace_original\": true}")
			}
		}

	case "button":
		var buttonRequestStruct RequestButtonStruct
		err := json.Unmarshal(ctx.Request.PostArgs().Peek("payload"), &buttonRequestStruct)
		if err != nil {
			fmt.Println(err)
		}
		switch {
		// If Team Admin Accepted the request
		case strings.HasPrefix(buttonRequestStruct.Actions[0].Value, "ACCEPTED_"):
			var (
				response_url = appRequest.ResponseURL
			)
			getAllDetails := strings.Split(string(buttonRequestStruct.Actions[0].Value), "_")[1]
			values := strings.Split(getAllDetails, ":")
			userId := values[0]
			teamName := values[1]
			roleType := values[2]
			teamId := git.GetTeamNamesAndIdsMap()[teamName]
			membership, err := git.AddTeamMembership(teamId, userId, roleType)
			if err != nil {
				client.HitRequest(response_url, "POST", header, "{\"text\" : \"Could not complete the request at this moment.\", \"response_type\": \"in_channel\", \"replace_original\": true}")
			} else {
				client.HitRequest(response_url, "POST", header, SubstParams([]string{
					userId,
					teamName,
					roleType,
					buttonRequestStruct.User.Name,
					"",
					values[3],
				}, GetPayload("changeAfterApproval.json")))
			}
			//Send status about the request back to user
			payload := SubstParams([]string{
				userId,
				teamName,
				roleType,
				buttonRequestStruct.User.Name,
				"Excellent!!! Your Request has been approved!!! Your membership status is :" + membership.GetState(),
				values[3],
			}, GetPayload("changeAfterApproval.json"))
			client.HitRequest(SLACK_WEBHOOK_TO_SEND_SLACKBOT, "POST", header, payload)

		// If Team Admin Declined the request
		case strings.HasPrefix(buttonRequestStruct.Actions[0].Value, "DECLINED_"):
			var (
				response_url = appRequest.ResponseURL
			)
			getAllDetails := strings.Split(string(buttonRequestStruct.Actions[0].Value), "_")[1]
			values := strings.Split(getAllDetails, ":")
			userId := values[0]
			teamName := values[1]
			roleType := values[2]

			fmt.Println(SubstParams([]string{
				userId,
				teamName,
				roleType,
				buttonRequestStruct.User.Name,
				"",
				values[3],
			}, GetPayload("changeAfterApproval.json")))

			//Change Existing Card
			client.HitRequest(response_url, "POST", header, SubstParams([]string{
				userId,
				teamName,
				roleType,
				buttonRequestStruct.User.Name,
				"",
				values[3],
			}, GetPayload("changeAfterDeclined.json")))

			//Send a notification to user.
			payload := SubstParams([]string{
				userId,
				teamName,
				roleType,
				buttonRequestStruct.User.Name,
				"Sorry!!! Your Request has been declined :( Please get in touch with Team Admin",
				values[3],
			}, GetPayload("changeAfterDeclined.json"))
			client.HitRequest(SLACK_WEBHOOK_TO_SEND_SLACKBOT, "POST", header, payload)

		case strings.HasPrefix(buttonRequestStruct.Actions[0].Value, "ACREATEREPO_"):
			var (
				response_url = appRequest.ResponseURL
			)
			getAllDetails := strings.Split(string(buttonRequestStruct.Actions[0].Value), "_")[1]
			values := strings.Split(getAllDetails, ":")

			name := values[0]
			description := values[1]
			private := values[2]
			teamName := values[3]
			teamId := git.GetTeamNamesAndIdsMap()[teamName]

			createrepoResult := git.CreateRepository(name, description, private, teamId)
			if !createrepoResult {
				client.HitRequest(response_url, "POST", header, "{\"text\" : \"Could not complete the request at this moment.\", \"response_type\": \"in_channel\", \"replace_original\": true}")
			} else {
				client.HitRequest(response_url, "POST", header, SubstParams([]string{
					name,
					private,
					teamName,
					buttonRequestStruct.User.Name,
					"",
					values[4],
				}, GetPayload("createRepoApproved.json")))
			}

			//Send status about the request back to user
			payload := SubstParams([]string{
				name,
				private,
				teamName,
				buttonRequestStruct.User.Name,
				"Excellent!!! Your Request has been approved!!!",
				values[4],
			}, GetPayload("createRepoApproved.json"))
			client.HitRequest(SLACK_WEBHOOK_TO_SEND_SLACKBOT, "POST", header, payload)

		case strings.HasPrefix(buttonRequestStruct.Actions[0].Value, "DCREATEREPO_"):
			var (
				response_url = appRequest.ResponseURL
			)
			getAllDetails := strings.Split(string(buttonRequestStruct.Actions[0].Value), "_")[1]
			values := strings.Split(getAllDetails, ":")

			name := values[0]
			private := values[2]
			teamName := values[3]

			//Change in githike channel
			client.HitRequest(response_url, "POST", header, SubstParams([]string{
				name,
				private,
				teamName,
				buttonRequestStruct.User.Name,
				"",
				values[4],
			}, GetPayload("createRepoApproved.json")))

			//Send status about the request back to user
			payload := SubstParams([]string{
				name,
				private,
				teamName,
				buttonRequestStruct.User.Name,
				"Sorry!!! Your Request has been declined :( Please get in touch with Team Admin",
				values[4],
			}, GetPayload("createRepoApproved.json"))
			client.HitRequest(SLACK_WEBHOOK_TO_SEND_SLACKBOT, "POST", header, payload)
		}
	}
}

func AddHikeTeamMembership(ctx *fasthttp.RequestCtx, db *sql.DB) {
	response_url := string(ctx.PostArgs().Peek("response_url"))
	textStr := fmt.Sprintf("%s", ctx.PostArgs().Peek("text"))
	commandLineParams := strings.Split(textStr, "#")
	if len(commandLineParams) != 3 {
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Commandline params are not valid :"+textStr+"`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
		return
	}
	teamName := commandLineParams[0]
	githubUserId := commandLineParams[1]
	roleType := commandLineParams[2]
	client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Checking if userid belongs to Hike...`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
	callerId := fmt.Sprintf("%s", ctx.PostArgs().Peek("user_id"))
	if ok, _ := git.CheckIfUserIsMemberOfOrg(commandLineParams[1]); ok {
		keysAndValues := git.GetTeamNamesAndIdsMap()
		if _, ok := keysAndValues[commandLineParams[0]]; !ok {
			client.HitRequest(response_url, "POST", header, "{ \"text\": \"`There in no such team in our config. Please check the Team's list again.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
			return
		}
		teamAdmin := GetTeamAdminFromDB(commandLineParams[0], db)
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Your request has been sent to : "+teamAdmin+"`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
		NotifyAdminAndUser(response_url, callerId, githubUserId, teamAdmin, teamName, roleType, "")
	} else {
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Can't perform this task as User does not belongs to Hike.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
		return
	}
}

func AddOrUpdateTeam(ctx *fasthttp.RequestCtx, db *sql.DB) {
	var (
		teamAdminStruct  TeamAdminMapStr
		teamAdminMapping = make(map[string]string)
		//Data            Response.&Data
	)
	type Data struct {
		RecordsUpdated int `json:"recordsUpdated"`
	}
	err := json.Unmarshal(ctx.Request.Body(), &teamAdminStruct)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range teamAdminStruct.Teamadminmap {
		if strings.HasSuffix(v.Admin, "@hike.in") {
			teamAdminMapping[v.Teamname] = v.Admin
		} else {
			SetErrorResponse(ctx, "3001", "ERROR", v.Admin+" Is not a HIKE user", http.StatusInternalServerError)
			return
		}
	}
	if res, records := DeleteAndAdd(teamAdminMapping, db); res {
		SetSuccessResponse(ctx, "1001", "SUCCESS", "Successfully inserted all records", http.StatusOK, &Data{
			RecordsUpdated: records,
		})
	} else {
		SetErrorResponse(ctx, "3001", "ERROR", "Error in inserting records", http.StatusInternalServerError)
	}
}

func DeleteMember(ctx *fasthttp.RequestCtx, db *sql.DB) {
	response_url := string(ctx.PostArgs().Peek("response_url"))
	textStr := fmt.Sprintf("%s", ctx.PostArgs().Peek("text"))
	callerId := fmt.Sprintf("%s", ctx.PostArgs().Peek("user_id"))

	//only admin can delete/deactivate github accounts
	if callerId == "U02A1MA8Z" || callerId == "U4XFTJW95" {
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Checking if userid belongs to Hike...`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
		//if email
		if strings.Contains(textStr, "@hike.in") {
			githubId := git.GetGithubIdFromEmail(strings.TrimSpace(textStr))
			fmt.Println("githubid from email --------->" + githubId)

			//If github id is empty
			if githubId == "" || githubId == "UNKNOWN" {
				client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Can't perform this task as no github id is associated with this email id. Use github id instead.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
				return
			}

			//Check is user belongs to hike
			if ok, _ := git.CheckIfUserIsMemberOfOrg(strings.TrimSpace(githubId)); !ok {
				client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Can't perform this task as User does not belongs to Hike.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
				return
			}

			//Finally deactivate the account
			if git.DeactivateGithubHikeAccount(githubId) {
				client.HitRequest(SLACK_WEBHOOK, "POST", header, "{ \"text\": \"`User Deactivated :"+githubId+"`\"}")
			}

		} else {
			//Check is user belongs to hike
			if ok, _ := git.CheckIfUserIsMemberOfOrg(strings.TrimSpace(textStr)); !ok {
				client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Can't perform this task as User does not belongs to Hike.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
				return
			}

			//Finally deactivate the account
			if git.DeactivateGithubHikeAccount(textStr) {
				client.HitRequest(SLACK_WEBHOOK, "POST", header, "{ \"text\": \"`User Deactivated :"+textStr+"`\"}")
			}
		}
		//if un-authorised
	} else {
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`You are not authorised to perform this task.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
	}
}

func CreateRepository(ctx *fasthttp.RequestCtx, db *sql.DB) {
	response_url := string(ctx.PostArgs().Peek("response_url"))
	textStr := fmt.Sprintf("%s", ctx.PostArgs().Peek("text"))
	commandLineParams := strings.Split(textStr, "#")
	if len(commandLineParams) != 4 {
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Commandline params are not valid :"+textStr+"`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
		return
	}
	name := commandLineParams[0]
	description := commandLineParams[1]
	private := commandLineParams[2]
	teamName := commandLineParams[3]
	callerId := fmt.Sprintf("%s", ctx.PostArgs().Peek("user_id"))
	keysAndValues := git.GetTeamNamesAndIdsMap()
	if _, ok := keysAndValues[commandLineParams[3]]; !ok {
		client.HitRequest(response_url, "POST", header, "{ \"text\": \"`There in no such team in our config. Please check the Team's list again.`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
		return
	}
	teamAdmin := GetTeamAdminFromDB(teamName, db)
	if teamAdmin == "" {
		teamAdmin = "abhishekg@hike.in"
	}
	client.HitRequest(response_url, "POST", header, "{ \"text\": \"`Your request has been sent to : "+teamAdmin+"`\", \"response_type\": \"ephemeral\", \"replace_original\": true }")
	//name, description, private, teamname, teamid
	NotifyAdminAndUserCreateRepoVersion(response_url, callerId, name, description, private, teamName, teamAdmin, keysAndValues[teamName])
}

//DeleteAndAdd ...
//Delete and add team admins
func DeleteAndAdd(teamadmin map[string]string, db *sql.DB) (bool, int) {
	query := "DELETE FROM teamadmins"
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error in deleting entries")
	}
	stmt, _ := db.Prepare("INSERT INTO teamadmins VALUES (null, ?, ?)")
	var counter int
	defer stmt.Close()
	for key := range teamadmin {
		res, _ := stmt.Exec(key, teamadmin[key])
		if result, _ := res.RowsAffected(); result > 0 {
			counter++
		}
	}
	if counter == len(teamadmin) {
		return true, counter
	} else {
		return false, counter
	}
}

func GetTeamAdminFromDB(teamname string, db *sql.DB) (admin string) {
	var (
		query = "SELECT admin FROM teamadmins WHERE team = '" + teamname + "'"
	)
	fmt.Println(query)
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
	}
	for rows.Next() {
		rows.Scan(&admin)
	}
	return
}

func ToJsonString(p interface{}) string {
	bytes, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error", err.Error())
	}
	return string(bytes)
}

func SearchInList(list []string, valueToBeSearched string) bool {
	var result bool = false
	for _, val := range list {
		fmt.Println(val, "    :   ", valueToBeSearched)
		if val == valueToBeSearched {
			result = true
		}
	}
	return result
}

func NotifyAdminAndUser(response_url, callerId, githubUserId, teamAdmin, teamName, roleType, message string) {
	//fmt.Println(response_url)
	options := make([]string, 4)
	options[0] = githubUserId + " wants approval of " + teamName + " as Role : " + roleType
	options[1] = "ACCEPTED_" + githubUserId + ":" + teamName + ":" + roleType + ":" + callerId
	options[2] = "DECLINED_" + githubUserId + ":" + teamName + ":" + roleType + ":" + callerId

	api := slack.New(SLACK_TOKEN)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range users {
		if teamAdmin == val.Profile.Email {
			options[3] = "Hey <@" + val.ID + ">,\n "
			payload := SubstParams(options, GetPayload("sendToAdmin.json"))
			client.HitRequest(SLACK_WEBHOOK, "POST", header, payload)
		}
	}
}

//name, description, private, teamname, teamid
func NotifyAdminAndUserCreateRepoVersion(response_url, callerId, name, description, private, teamname, teamAdmin string, teamid int) {
	//fmt.Println(response_url)
	options := make([]string, 4)
	options[0] = GetEmailIdFromSlackId(callerId) + " wants to create *" + name + "* repository under  : *" + teamname + "*"
	options[1] = "ACREATEREPO_" + name + ":" + description + ":" + private + ":" + teamname + ":" + callerId
	options[2] = "DCREATEREPO_" + name + ":" + description + ":" + private + ":" + teamname + ":" + callerId

	api := slack.New(SLACK_TOKEN)
	users, err := api.GetUsers()
	if err != nil {
		fmt.Println(err)
	}
	for _, val := range users {
		if teamAdmin == val.Profile.Email {
			options[3] = "Hey <@" + val.ID + ">,\n "
			payload := SubstParams(options, GetPayload("sendToAdmin.json"))
			client.HitRequest(SLACK_WEBHOOK, "POST", header, payload)
		}
	}
}
