package sender

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func SendMessage(s *discordgo.Session, channelID string, message string) {
	if len(message) >= 2000 {
		sendMessageInChunks(s, channelID, message)
	} else {
		_, err := s.ChannelMessageSend(channelID, message)
		HandleErr(s, channelID, err)
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

func SendReply(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	if len(message) >= 2000 {
		sendReplyInChunks(s, m, message)
	} else {
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		HandleErr(s, m.ChannelID, err)
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

func HandleErr(s *discordgo.Session, channelID string, err error) bool {
	if err != nil {
		log.Printf("Received error: %s", err.Error())
		s.ChannelMessageSend(channelID, "Error:\n```"+err.Error()+"```")
		return true
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
