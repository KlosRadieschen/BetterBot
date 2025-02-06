package messages

import (
	ai "BetterScorch/ai"
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var responses []messageResponse

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
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
		handleErr(s, m, err)
		sendPotentiallyBigAssMessage(s, m, resp)
	}
}

func sendPotentiallyBigAssMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if len(message) >= 2000 {
		sendInChunks(s, m, message)
	} else {
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		handleErr(s, m, err)
	}
}

func sendInChunks(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	chunks := splitIntoChunks(message, 1999)

	for _, chunk := range chunks {
		s.ChannelMessageSendReply(m.ChannelID, chunk, m.Reference())
	}
}

func splitIntoChunks(message string, chunkSize int) []string {
	var chunks []string
	for len(message) > 0 {
		end := min(len(message), chunkSize)
		chunks = append(chunks, message[:end])
		message = message[end:]
	}
	return chunks
}

func handleErr(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Error:\n```"+err.Error()+"```", m.Reference())
		return
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
