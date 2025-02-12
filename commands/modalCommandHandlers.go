package commands

import (
	"BetterScorch/polls"
	"BetterScorch/sender"
	"time"

	"github.com/bwmarrin/discordgo"
)

func inputPollModalCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Poll",
			Content:  i.Message.Content,
			CustomID: "inputpollmodalsubmit",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "response",
							Label:    "Response",
							Style:    discordgo.TextInputShort,
							Required: true,
						},
					},
				},
			},
		},
	})
	sender.HandleErrInteraction(s, i, err)
}

func inputPollModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	polls.SubmitInputPollResponse(s, i.Message.ID, i.Member.User.ID, response)
	sender.RespondEphemeral(s, i, "Answer submitted")
	sender.SetResponseTimeout(s, i, 5*time.Second)
}

func inputPollShowHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.ThinkEphemeral(s, i)
	embeds, err := polls.GetAllInputsEmbeds(s, i.Message.ID)
	if err != nil {
		sender.RespondEphemeral(s, i, "Sorry, his poll is broken")
	} else {
		if polls.GetNumberOfInputs(i.Message.ID) == 0 {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: "No responses",
			})
			sender.SetResponseTimeout(s, i, 3*time.Second)
		} else {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Flags:  discordgo.MessageFlagsEphemeral,
				Embeds: embeds,
			})
		}
	}
}
