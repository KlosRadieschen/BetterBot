package messages

import (
	"BetterScorch/ai"
	"BetterScorch/execution"
	"BetterScorch/sender"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if execution.CheckAndDeleteExecuteeMessage(s, m) {
		return
	}
	handleAIResponses(s, m)

	for _, response := range responses {
		for _, trigger := range response.triggers {
			// Magical RegEx bullshiterry
			match := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(trigger))).FindStringSubmatch(strings.ToLower(m.Content))
			if match != nil {
				response.handleResponse(s, m)
			}
		}
	}
}

func handleAIResponses(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Type == 19 && m.ReferencedMessage.Author.ID == "1196526025211904110" || strings.Contains(strings.ToLower(m.Content), "scorch") || strings.Contains(strings.ToLower(m.Content), "dementia") {
		resp, err := ai.GenerateResponse(m.Content)
		sender.HandleErr(s, m.ChannelID, err)
		sender.SendReply(s, m, resp)
	}
}
