package sender

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

var webhook *discordgo.Webhook

func InitWebhook(s *discordgo.Session) {
	webhook, _ = s.WebhookCreate("1195135473643958316", "Scorch", "https://media.sketchfab.com/models/be06a067516e4084bda252f2bb8ed008/thumbnails/485c890fd0354081b28001171695ecaa/22a61b89d3d44964be605a360e77c4c3.jpeg")
}

func SendPersonalityMessage(s *discordgo.Session, channelID string, msg string, name string, pfpLink string, req *openai.ChatCompletionRequest) {
	_, err := s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Username:  name,
		Content:   msg,
		AvatarURL: pfpLink,
	})
	HandleErr(s, channelID, err)
}

func SendPersonalityReply(s *discordgo.Session, m *discordgo.MessageCreate, msg string, name string, pfpLink string, req *openai.ChatCompletionRequest) {
	_, err := s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("> [Replying to](https://discord.com/channels/@me/%v/%v): %v\n\n%v",
			m.ChannelID,
			m.ID,
			m.Author.Mention(),
			msg,
		),
		Username:  name,
		AvatarURL: pfpLink,
	})
	HandleErr(s, m.ChannelID, err)
}
