package sender

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

func SendPersonalityMessage(s *discordgo.Session, channelID string, msg string, name string, pfpLink string, req *openai.ChatCompletionRequest) {
	wh, err := s.WebhookCreate(channelID, name, pfpLink)
	if HandleErr(s, channelID, err) {
		return
	}
	_, err = s.WebhookExecute(wh.ID, wh.Token, false, &discordgo.WebhookParams{
		Content: msg,
	})
	HandleErr(s, channelID, err)
	err = s.WebhookDelete(wh.ID)
	HandleErr(s, channelID, err)
}

func SendPersonalityReply(s *discordgo.Session, m *discordgo.MessageCreate, msg string, name string, pfpLink string, req *openai.ChatCompletionRequest) {
	wh, err := s.WebhookCreate(m.ChannelID, name, pfpLink)
	if HandleErr(s, m.ChannelID, err) {
		return
	}

	_, err = s.WebhookExecute(wh.ID, wh.Token, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("> [Replying to](https://discord.com/channels/@me/%v/%v): %v\n\n%v",
			m.ChannelID,
			m.ID,
			m.Author.Mention(),
			msg,
		),
		AvatarURL: pfpLink,
	})
	HandleErr(s, m.ChannelID, err)
	err = s.WebhookDelete(wh.ID)
	HandleErr(s, m.ChannelID, err)
}
