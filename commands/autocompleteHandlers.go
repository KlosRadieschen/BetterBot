package commands

import (
	"BetterScorch/webhooks"

	"github.com/bwmarrin/discordgo"
)

func messageAutocompleteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for ID, characters := range webhooks.CharacterBuffer {
		if ID == i.Member.User.ID {
			choices := []*discordgo.ApplicationCommandOptionChoice{}
			for _, character := range characters {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  character.Name,
					Value: character.Name,
				})
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionApplicationCommandAutocompleteResult,
				Data: &discordgo.InteractionResponseData{
					Choices: choices,
				},
			})

			if err != nil {
				panic(err)
			}
		}
	}
}
