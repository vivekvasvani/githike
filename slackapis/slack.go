package slackapis

import (
	"fmt"

	"github.com/nlopes/slack"
)

const (
	SLACK_TOKEN = "xoxp-2151902985-167537642311-181450648469-1abc20fff346031bb10664d62f64e997"
	OPTIONS     = "I support the following commands:\n" +
		"`!DeleteKey` - Delete Deploy Key from a repo.\n" +
		"`!AddKey` - Add Deploy Key for a repo.\n" +
		"`!CreateRepo` - Creates a new PRIVATE [default] repository in the specified GitHub organization.\n" +
		"`!ListPRs` - List the Pull Requests for a repo.\n" +
		"`!ListOrgs` - Lists the GitHub organizations that are managed.\n" +
		"`!SetHomepage` - Adds/Modifies a GitHub repo's homepage URL.\n" +
		"`!SetDescription` - Adds/Modifies a GitHub repo's description.\n" +
		"`!SetBranchProtection` - Toggles the branch protection for a repo.\n" +
		"`!AddUserToTeam` - Adds a GitHub user to a specific team inside the organization.\n" +
		"`!AddCollab` - Adds an outside collaborator to a specific repository in a specific GitHub organization.\n" +
		"`!ListKeys` - List the Deploy Keys for a repo.\n" +
		"`!SetDefaultBranch` - Sets the default branch for a repo.\n" +
		"`!GetKey` - Get Deploy Key Public Key\n" +
		"`!Help` - This command.\n"
)

func InviteUserToChannel() { //channelName, userID string) {
	api := slack.New(SLACK_TOKEN)
	channels, err := api.GetChannels(false)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	for _, channel := range channels {
		//fmt.Println(channel.ID, channel.Name)
		if channel.Name == "_githike_" {
			fmt.Println("Yesssss")
		}
	}
	channel, err := api.InviteUserToChannel("G5K9BA18E", "U4ZUFVBEW")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(channel.Members)

}
