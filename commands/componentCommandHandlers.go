package commands

import (
	"BetterScorch/database"
	"BetterScorch/polls"
	"BetterScorch/sender"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func pollVoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	polls.PollMutex.Lock()

	numString := i.MessageComponentData().CustomID[len(i.MessageComponentData().CustomID)-1]
	num, _ := strconv.Atoi(string(numString))
	err := polls.AddVote(s, i.Message.ID, i.Member.User.ID, num)
	if err != nil {
		sender.FollowupEphemeral(s, i, "Sorry, this poll is broken")
		go sender.SetResponseTimeout(s, i, 3*time.Second)
	} else {
		sender.FollowupEphemeral(s, i, "Votes updated")
		go sender.SetResponseTimeout(s, i, 3*time.Second)
	}
	polls.PollMutex.Unlock()
}

func inputPollVoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	polls.PollMutex.Lock()

	err := polls.SubmitInputPollResponse(s, i.Message.ID, i.Member.User.ID, i.MessageComponentData().Values[0])
	if err != nil {
		sender.FollowupEphemeral(s, i, "Sorry, this poll is broken")
		go sender.SetResponseTimeout(s, i, 3*time.Second)
	} else {
		sender.FollowupEphemeral(s, i, "Response submitted")
		go sender.SetResponseTimeout(s, i, 3*time.Second)
	}
	polls.PollMutex.Unlock()
}

func pollShowHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	polls.PollMutex.Lock()

	embeds, err := polls.GetAllVotesEmbeds(s, i.Message.ID)
	if err != nil {
		sender.RespondEphemeral(s, i, "Sorry, his poll is broken", nil)
	} else {
		if polls.GetVotesSum(i.Message.ID) == 0 {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "No responses",
			})
			go sender.SetResponseTimeout(s, i, 3*time.Second)
		} else {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: embeds,
			})
		}
	}
	polls.PollMutex.Unlock()
}

func inputPollShowHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	polls.PollMutex.Lock()

	embeds, err := polls.GetAllInputsEmbeds(s, i.Message.ID)
	if err != nil {
		sender.FollowupEphemeral(s, i, "Sorry, his poll is broken")
		go sender.SetResponseTimeout(s, i, 3*time.Second)
	} else {
		if polls.GetNumberOfInputs(i.Message.ID) == 0 {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "No responses",
			})
			go sender.SetResponseTimeout(s, i, 3*time.Second)
		} else {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: embeds,
			})
		}
	}
	polls.PollMutex.Unlock()
}

func decryptHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !isSWAG(i.Member) && !isHC(i.Member) {
		sender.RespondEphemeral(s, i, "https://tenor.com/view/yellow-emoji-no-no-emotiguy-no-no-no-gif-gif-9742000569423889376", nil)
		return
	}

	ID, _ := strings.CutPrefix(i.MessageComponentData().CustomID, "decrypt")
	r, _ := database.Get("EncryptedMessage", []string{"msg"}, &database.DBValue{Name: "msgID", Value: ID})

	embed := &discordgo.MessageEmbed{
		Title:       "Decrypted Message",
		Color:       16738740,
		Description: r[0][0],
	}

	sender.RespondEphemeral(s, i, "", []*discordgo.MessageEmbed{embed})
}
