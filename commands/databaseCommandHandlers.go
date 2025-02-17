package commands

import (
	"BetterScorch/sender"
	"BetterScorch/webhooks"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func addCharacterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var pfpLink string
	for _, attachment := range i.ApplicationCommandData().Resolved.Attachments {
		pfpLink = attachment.URL
	}

	if !strings.Contains(i.ApplicationCommandData().Options[1].StringValue(), "text") {
		sender.RespondError(s, i, "The user wanted to add a new character, but the \"brackets\" did not include the word text. For example, if you have a character named \"Scorch\", your brackets could be \"Scorch: text\", in which case a message with that pattern will be sent using Scorch (text will be the actual text of the message)")
		return
	}

	err := webhooks.AddCharacter(i.Member.User.ID, webhooks.Character{
		OwnerID:    i.Member.User.ID,
		Name:       i.ApplicationCommandData().Options[0].StringValue(),
		Brackets:   i.ApplicationCommandData().Options[1].StringValue(),
		AvatarLink: pfpLink,
	})
	if !sender.HandleErrInteraction(s, i, err) {
		sender.Respond(s, i, "Character created")
	}
}

func removeCharacterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if err := webhooks.RemoveCharacter(i.Member.User.ID, i.ApplicationCommandData().Options[0].StringValue()); err != nil {
		sender.HandleErrInteraction(s, i, err)
		return
	}

	sender.Respond(s, i, "Character removed")
}

func listCharactersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	characters, err := webhooks.ListCharacters(i.Member.User.ID)
	if err != nil {
		sender.HandleErrInteraction(s, i, err)
		return
	}

	if len(characters) == 0 {
		sender.Respond(s, i, "You have no characters")
		return
	}

	var message string
	for _, character := range characters {
		message += fmt.Sprintf("%s: %s\n", character.Name, character.Brackets)
	}

	sender.Respond(s, i, message)
}
