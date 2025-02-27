package commands

import (
	"fmt"
	"strconv"
	"strings"

	"BetterScorch/database"
	"BetterScorch/sender"
	"BetterScorch/webhooks"

	"github.com/bwmarrin/discordgo"
)

/*

Characters

*/

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

	sender.RespondEphemeral(s, i, "", embeds)
}

/*

Register

*/

func registerHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var pfpLink string
	for _, attachment := range i.ApplicationCommandData().Resolved.Attachments {
		pfpLink = attachment.URL
	}

	err := database.Insert("Pilot", []*database.DBValue{
		{Name: "name", Value: i.ApplicationCommandData().Options[0].StringValue()},
		{Name: "callsign", Value: i.ApplicationCommandData().Options[1].StringValue()},
		{Name: "picture", Value: pfpLink},
	}...)

	if err != nil {
		sender.HandleErrInteraction(s, i, err)
	} else {
		sender.Respond(s, i, " Successfully registered", nil)
	}
}

func unregisterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	affected, err := database.Remove("Pilot", []*database.DBValue{
		{Name: "pk_userID", Value: i.Member.User.ID},
	}...)

	if err != nil {
		sender.HandleErrInteraction(s, i, err)
	} else if affected == 0 {
		sender.RespondError(s, i, "User tried to delete himself from the database but never even registered")
	} else {
		sender.Respond(s, i, "Removed you from the database", nil)
	}
}

/*

REPORTS

*/

func listReportsHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	reports, err := database.GetAll("Report")
	if !sender.HandleErrInteraction(s, i, err) {
		if len(reports) == 0 {
			sender.Respond(s, i, "No results", nil)
			return
		}

		var resultString string
		for _, row := range reports {
			name := row[0]
			timeIndex, _ := strconv.Atoi(row[1])
			reportType := row[2]
			authorType := row[3]

			timeString := fmt.Sprintf("%v%v", ifNegative(timeIndex), timeIndex)
			resultString += fmt.Sprintf("- #%v%v%v: %v\n", reportType, timeString, authorType, name)
		}

		sendResponse(s, i, resultString)
	}
}

func getReportHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	numberString := fmt.Sprintf("%v", i.ApplicationCommandData().Options[0].IntValue())
	timeIndex := strings.TrimSuffix(strings.TrimPrefix(numberString, string(numberString[0])), string(numberString[len(numberString)-1]))
	var timeInt int
	fmt.Println(timeIndex)

	timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "0"))
	timeInt, _ = strconv.Atoi(strings.TrimPrefix(timeIndex, "1"))
	if timeIndex[0] == '0' {
		timeInt = -timeInt
	}

	fmt.Println(timeIndex)

	reports, err := database.Get("Report", []string{"pk_name", "fk_pilot_wrote", "description"}, database.DBValue{Name: "timeIndex", Value: strconv.Itoa(timeInt)})

	if !sender.HandleErrInteraction(s, i, err) {
		if len(reports) == 0 {
			sender.RespondError(s, i, "No results")
			return
		}

		row := reports[0]
		name := row[0]
		id := row[1]
		description := row[2]

		member, _ := s.State.Member(i.GuildID, id)
		nick := "Probably Saturn"
		if member != nil && member.Nick != "" {
			nick = member.Nick
		}

		embed := &discordgo.MessageEmbed{
			Title: name,
			Author: &discordgo.MessageEmbedAuthor{
				Name:    nick,
				URL:     fmt.Sprintf("https://aha-rp.org/get/pilots/%v", strings.ReplaceAll(nick, " ", "")),
				IconURL: fmt.Sprintf("https://aha-rp.org/static/assets/avatars/%v.png", id),
			},
			Color:       16738740,
			URL:         fmt.Sprintf("https://aha-rp.org/get/reports/%v", i.ApplicationCommandData().Options[0].IntValue()),
			Description: description,
		}

		if len(embed.Description) > 4000 {
			sender.Respond(s, i, fmt.Sprintf("Fuck you, the report is too long. Go read it here: https://aha-rp.org/get/reports/%v", i.ApplicationCommandData().Options[0].IntValue()), nil)
		} else {
			sender.Respond(s, i, "", []*discordgo.MessageEmbed{embed})
		}
	}
}

func ifNegative(index int) string {
	if index < 0 {
		return "0"
	}
	return "1"
}

func sendResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	if len(content) >= 2000 {
		chunks := make([]string, 0, len(content)/2000+1)
		currentChunk := ""
		for _, c := range content {
			if len(currentChunk) >= 1999 {
				chunks = append(chunks, currentChunk)
				currentChunk = ""
			}
			currentChunk += string(c)
		}
		if currentChunk != "" {
			chunks = append(chunks, currentChunk)
		}

		sender.Respond(s, i, chunks[0], nil)

		for _, chunk := range chunks[1:] {
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Content: chunk,
			})
		}
	} else {
		sender.Respond(s, i, content, nil)
	}
}
