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
