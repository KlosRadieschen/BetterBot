package commands

import (
	"BetterScorch/polls"
	"BetterScorch/sender"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

func pollVoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	numString := i.MessageComponentData().CustomID[len(i.MessageComponentData().CustomID)-1]
	num, _ := strconv.Atoi(string(numString))
	err := polls.AddVote(s, i.Message.ID, i.Member.User.ID, num)
	if err != nil {
		sender.FollowupEphemeral(s, i, "Sorry, this poll is broken")
		sender.SetResponseTimeout(s, i, 3*time.Second)
	} else {
		sender.FollowupEphemeral(s, i, "Votes updated")
		sender.SetResponseTimeout(s, i, 3*time.Second)
	}
}

func inputPollVoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)

	err := polls.SubmitInputPollResponse(s, i.Message.ID, i.Member.User.ID, i.MessageComponentData().Values[0])
	if err != nil {
		sender.FollowupEphemeral(s, i, "Sorry, this poll is broken")
		sender.SetResponseTimeout(s, i, 3*time.Second)
	} else {
		sender.FollowupEphemeral(s, i, "Response submitted")
		sender.SetResponseTimeout(s, i, 3*time.Second)
	}
}

func pollShowHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	votes, err := polls.GetAllVotesString(i.Message.ID)
	if err != nil {
		sender.RespondEphemeral(s, i, "Sorry, his poll is broken")
	} else {
		sender.RespondEphemeral(s, i, votes)
	}
}
