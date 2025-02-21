package commands

import (
	"BetterScorch/database"
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
		sender.RespondError(s, i, "The user wanted to add a new character, but the \"brackets\" did not include the word text. For example, if you have a character named \"Scorch\", your brackets COULD be \"Scorch: text\" (it doesn't matter where the word text is as long as it is there), which means that if the user writes the message \"Scorch: example\" the user's message will be replaced by the character's message \"example\". Mention the entire example. This response can be longer (more than a paragraph), so you can make sure that you explain it thoroughly and comprehensively.")
		return
	}

	err := webhooks.AddCharacter(i.Member.User.ID, webhooks.Character{
		OwnerID:    i.Member.User.ID,
		Name:       i.ApplicationCommandData().Options[0].StringValue(),
		Brackets:   i.ApplicationCommandData().Options[1].StringValue(),
		AvatarLink: pfpLink,
	})
	if !sender.HandleErrInteraction(s, i, err) {
		sender.Respond(s, i, "Character created", nil)
	}
}

func removeCharacterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if exists, err := webhooks.RemoveCharacter(i.Member.User.ID, i.ApplicationCommandData().Options[0].StringValue()); err != nil {
		sender.HandleErrInteraction(s, i, err)
	} else if !exists {
		sender.RespondError(s, i, "User tried to delete a character that does not exist")
	} else {
		sender.Respond(s, i, "Character removed", nil)
	}
}

func listCharactersHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	characters := webhooks.ListCharacters(i.Member.User.ID)

	if len(characters) == 0 {
		sender.RespondError(s, i, "User tried to list his characters but has none")
		return
	}

	embeds := []*discordgo.MessageEmbed{}

	for _, character := range characters {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title:       character.Name,
			Description: fmt.Sprintf("Brackets: \"%s\"", character.Brackets),
			Image:       &discordgo.MessageEmbedImage{URL: character.AvatarLink},
		})
	}

	sender.Respond(s, i, "", embeds)
}

func registerHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var pfpLink string
	for _, attachment := range i.ApplicationCommandData().Resolved.Attachments {
		pfpLink = attachment.URL
	}

	err := database.Insert("Pilot", []*database.DBValue{
		{
			Name:  "name",
			Value: i.ApplicationCommandData().Options[0].StringValue(),
		},
		{
			Name:  "callsign",
			Value: i.ApplicationCommandData().Options[1].StringValue(),
		},
		{
			Name:  "picture",
			Value: pfpLink,
		},
	}...)

	if err != nil {
		sender.HandleErrInteraction(s, i, err)
	} else {
		sender.Respond(s, i, " Successfully registered", nil)
	}
}

func unregisterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	affected, err := database.Remove("Pilot", []*database.DBValue{
		{
			Name:  "pk_userID",
			Value: i.Member.User.ID,
		},
	}...)

	if err != nil {
		sender.HandleErrInteraction(s, i, err)
	} else if affected == 0 {
		sender.RespondError(s, i, "User tried to delete himself from the database but never even registered")
	} else {
		sender.Respond(s, i, "Removed you from the database", nil)
	}
}
