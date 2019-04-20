package utils

import (
	"github.com/bwmarrin/discordgo"
)

var BotID string = "556662039079157772"

var AllBotIDs []string = []string{"279722369260453888", BotID, "481910793697493011"}

func UserIsGod(userid string) bool {
	return userid == "190192176632692740"
}

func UserIsBot(userid string) bool {
	return userid == BotID
}

func ChannelIsTest(channelid string) bool {
	return channelid == "556672189265477662"
}

func FilterBotsFromMembers(members []*discordgo.Member) []*discordgo.Member {
	newMembers := []*discordgo.Member{}
	for _, member := range members {
		if !Contains(AllBotIDs, member.User.ID) {
			newMembers = append(newMembers, member)
		}
	}
	return newMembers
}
