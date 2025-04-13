package commands

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"BetterScorch/database"
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"BetterScorch/stocks"
	"BetterScorch/webhooks"

	"github.com/bwmarrin/discordgo"
)

/*

Characters

*/

func addCharacterHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var pfpLink string
	for _, attachment := range i.ApplicationCommandData().Resolved.Attachments {
		// Fetch the file from the URL
		resp, err := http.Get(attachment.URL)
		if sender.HandleErrInteraction(s, i, err) {
			return
		}
		defer resp.Body.Close()

		// Read the file into memory
		fileData, err := io.ReadAll(resp.Body)
		if sender.HandleErrInteraction(s, i, err) {
			return
		}

		// Create a file to send
		file := &discordgo.File{
			Name:   attachment.Filename,
			Reader: bytes.NewReader(fileData),
		}

		// Send the file
		msg, _ := s.ChannelMessageSendComplex("1196943729387372634", &discordgo.MessageSend{
			Files: []*discordgo.File{file},
		})

		pfpLink = msg.Attachments[0].URL
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

Gambling

*/

func entereconomyHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := stocks.Enter(i.Member.User.ID)

	if err != nil {
		sender.RespondError(s, i, "The user is trying to enter the ScorchCoin economy but is already in it")
	} else {
		sender.Respond(s, i, "Successfully entered the ScorchCoin economy with 420 ScorchCoin", nil)
	}
}

func balanceHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	balance, err := stocks.ModifyBalance(i.Member.User.ID, 0)
	if !sender.HandleErrInteraction(s, i, err) {
		sender.Respond(s, i, fmt.Sprintf("You currently own %v ScorchCoin", balance), nil)
	}
}

func stonksHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// You cant use options with boolean values so we parse an int instead
	bool := false
	if i.ApplicationCommandData().Options[2].IntValue() == 1 {
		bool = true
	}

	if bool {
		_, err := stocks.ModifyBalance(i.Member.User.ID, -int(i.ApplicationCommandData().Options[1].IntValue()))
		if err != nil {
			sender.RespondError(s, i, err.Error())
		}
	} else {
		_, err := stocks.ModifyBalance(i.Member.User.ID, int(i.ApplicationCommandData().Options[1].IntValue()))
		if err != nil {
			sender.RespondError(s, i, err.Error())
		}
	}

	newVal, err := stocks.Trade(i.Member.User.ID, i.ApplicationCommandData().Options[0].StringValue(), int(i.ApplicationCommandData().Options[1].IntValue()), bool)
	if err != nil {
		sender.RespondError(s, i, err.Error())
	}

	sender.Respond(s, i, fmt.Sprintf("You now own %v of Scorchcoin in %v", newVal, i.ApplicationCommandData().Options[0].StringValue()), nil)
}

func gambleHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Options[1].IntValue() {
	case 0:
		coinflipHandler(s, i)
	}
}

func coinflipHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Respond(s, i, "https://tenor.com/view/eminem-eminem-taern-eminem-taern-coin-toss-eminem-imperial-imperial-gif-13493891068666315081", nil)
	time.Sleep(3 * time.Second)
	if rand.Intn(2) == 0 {
		newVal, err := stocks.ModifyBalance(i.Member.User.ID, -int(i.ApplicationCommandData().Options[0].IntValue()))
		if !sender.HandleErrInteractionFollowup(s, i, err) {
			sender.Followup(s, i, fmt.Sprintf("YOU LOSE %v SCORCHCOIN (New Balance: %v)", i.ApplicationCommandData().Options[0].IntValue(), newVal))
		}
	} else {
		newVal, err := stocks.ModifyBalance(i.Member.User.ID, int(i.ApplicationCommandData().Options[0].IntValue()))
		if !sender.HandleErrInteractionFollowup(s, i, err) {
			sender.Followup(s, i, fmt.Sprintf("YOU WIN %v SCORCHCOIN (New Balance: %v)", i.ApplicationCommandData().Options[0].IntValue(), newVal))
		}
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

	reports, err := database.Get("Report", []string{"pk_name", "fk_pilot_wrote", "description"}, &database.DBValue{Name: "timeIndex", Value: strconv.Itoa(timeInt)})

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

// Legacy
func addReportHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v:3306)/%v", secrets.DBUser, secrets.DBPassword, secrets.DBAddress, secrets.DBName)
	var err error
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert data into the table
	stmt, err := db.Prepare("SELECT MAX(timeIndex) FROM Report")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	// Execute the prepared statement with actual values
	var maxIndex int
	rows, err := stmt.Query()
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&maxIndex); err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: err.Error(),
				},
			})
			return
		}
	}

	maxIndex += 10

	member, _ := s.GuildMember("1195135473006420048", i.Member.User.ID)
	var roles []string
	var authorIndex int
	index := -1
	roles = append(roles, "1195135956471255140")
	roles = append(roles, "1195858311627669524")
	roles = append(roles, "1195858271349784639")
	roles = append(roles, "1195136106811887718")
	roles = append(roles, "1195858179590987866")
	roles = append(roles, "1195137362259349504")
	roles = append(roles, "1195136284478410926")
	roles = append(roles, "1195137253408768040")
	roles = append(roles, "1195758308519325716")
	roles = append(roles, "1195758241221722232")
	roles = append(roles, "1195758137563689070")
	roles = append(roles, "1195757362439528549")
	roles = append(roles, "1195136491148550246")
	roles = append(roles, "1195708423229165578")
	roles = append(roles, "1195137477497868458")
	roles = append(roles, "1195136604373782658")
	roles = append(roles, "1195711869378367580")

	for i, guildRole := range roles {
		for _, memberRole := range member.Roles {
			if guildRole == memberRole {
				index = i
			}
		}
	}

	if index >= 0 && index <= 3 {
		authorIndex = 1
	} else if index >= 4 && index <= 7 {
		authorIndex = 2
	} else if index >= 8 && index <= 11 {
		authorIndex = 3
	} else if index >= 12 && index <= 14 {
		authorIndex = 4
	} else {
		authorIndex = 5
	}

	stmt, err = db.Prepare("INSERT INTO Report VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	// Execute the prepared statement with actual values
	name := i.ApplicationCommandData().Options[0].StringValue()
	reportType := i.ApplicationCommandData().Options[1].IntValue()
	report := i.ApplicationCommandData().Options[2].StringValue()

	if reportType >= 10 || reportType < 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Listen up, you insolent excuse for a pilot! You dare insult me, the mighty Scorch, and then have the audacity to try adding a report with an invalid 'type' number? Are you malfunctioning or just plain stupid? Let me spell it out for you since you seem to be lacking basic cognitive functions: the 'type' number is only ONE DIGIT! How hard is it to understand that?! If you can't even get that simple detail right, I shudder to think about your piloting skills. Fix your mistake immediately before I decide to unleash my fury upon you and your sorry excuse for a Titan! Now, get it together, or face the consequences!",
			},
		})
		return
	}

	_, err = stmt.Exec(&name, &maxIndex, &reportType, &authorIndex, &i.Member.User.ID, &report)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: err.Error(),
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Report added",
		},
	})
}
