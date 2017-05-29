package slackapis

import (
	"fmt"

	"github.com/nlopes/slack"
)

const (
	SLACK_TOKEN = ""
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
