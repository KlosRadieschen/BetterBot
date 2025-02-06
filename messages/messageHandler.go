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
		chunks := make([]string, 0, len(message)/2000+1)
		currentChunk := ""
		for _, c := range message {
			if len(currentChunk) >= 1999 {
				chunks = append(chunks, currentChunk)
				currentChunk = ""
			}
			currentChunk += string(c)
		}
		if currentChunk != "" {
			chunks = append(chunks, currentChunk)
		}
		for _, chunk := range chunks[0:] {
			s.ChannelMessageSendReply(m.ChannelID, chunk, m.Reference())
		}
	} else {
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func handleErr(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	if err != nil {
		s.ChannelMessageSendReply(m.ChannelID, "Error:\n```"+err.Error()+"```", m.Reference())
		return
	}
}
