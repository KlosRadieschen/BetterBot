package commands

import (
	"BetterScorch/messages"
	"BetterScorch/polls"
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"BetterScorch/webhooks"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func inputPollModalCreateHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := polls.GetAllInputsEmbeds(s, i.Message.ID)
	if err != nil {
		sender.RespondEphemeral(s, i, "Sorry, his poll is broken", nil)
		go sender.SetResponseTimeout(s, i, 3*time.Second)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
	sender.RespondEphemeral(s, i, "Answer submitted", nil)
	go sender.SetResponseTimeout(s, i, 5*time.Second)
}

func messageModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	split := strings.Split(i.ModalSubmitData().CustomID, "-")
	recipientID := split[1]
	message := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	var name string
	var pfp string
	if len(split) > 2 {
		for ID, characters := range webhooks.CharacterBuffer {
			if ID == i.Member.User.ID {
				for _, character := range characters {
					if character.Name == split[2] {
						name = character.Name
						pfp = character.AvatarLink
					}
				}
			}
		}
	} else {
		member, _ := s.GuildMember(secrets.GuildID, i.Member.User.ID)
		name = member.Nick
		pfp = member.AvatarURL("")
	}

	messages.UserMessages = append(messages.UserMessages, messages.UserMessage{
		SenderName:  name,
		SenderPFP:   pfp,
		RecipientID: recipientID,
		Message:     message,
	})

	sender.RespondEphemeral(s, i, "Message send", nil)
	go sender.SetResponseTimeout(s, i, 5*time.Second)
}
