package commands

import (
	"BetterScorch/polls"
	"BetterScorch/sender"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func pollVoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	numString := i.MessageComponentData().CustomID[len(i.MessageComponentData().CustomID)-1]
	num, _ := strconv.Atoi(string(numString))
	err := polls.AddVote(s, i.Message.ID, i.Member.User.ID, num)
	if err != nil {
		if strings.Contains(err.Error(), "Repeated vote") {
			sender.FollowupEphemeral(s, i, "You have already voted on this poll, idiot")
		} else {
			sender.FollowupEphemeral(s, i, "This poll is broken, sorry")
		}
	} else {
		sender.FollowupEphemeral(s, i, "Votes updated")
	}
}
