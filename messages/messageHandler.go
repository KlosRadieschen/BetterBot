package messages

import (
	"BetterScorch/ai"
	"BetterScorch/execution"
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
		handleErr(s, m.ChannelID, err)
		sendReply(s, m, resp)
	}
}

func SendMessage(s *discordgo.Session, channelID string, message string) {
	if len(message) >= 2000 {
		sendMessageInChunks(s, channelID, message)
	} else {
		_, err := s.ChannelMessageSend(channelID, message)
		handleErr(s, channelID, err)
	}
}

func sendMessageInChunks(s *discordgo.Session, channelID string, message string) {
	chunks := splitIntoChunks(message, 1999)
	msg, _ := s.ChannelMessageSend(channelID, chunks[0])
	previous := msg.Reference()

	for i, chunk := range chunks {
		if i == 0 {
			break
		}
		msg, _ := s.ChannelMessageSendReply(channelID, chunk, previous)
		previous = msg.Reference()
	}
}

func sendReply(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if len(message) >= 2000 {
		sendReplyInChunks(s, m, message)
	} else {
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		handleErr(s, m.ChannelID, err)
	}
}

func sendReplyInChunks(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	chunks := splitIntoChunks(message, 1999)
	previous := m.Reference()

	for _, chunk := range chunks {
		msg, _ := s.ChannelMessageSendReply(m.ChannelID, chunk, previous)
		previous = msg.Reference()
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

func handleErr(s *discordgo.Session, channelID string, err error) {
	if err != nil {
		s.ChannelMessageSend(channelID, "Error:\n```"+err.Error()+"```")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
