package sender

import (
	"BetterScorch/secrets"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sashabaranov/go-openai"
)

var webhook *discordgo.Webhook

func InitWebhook(s *discordgo.Session) {
	var err error
	webhook, err = s.Webhook("1339582520110481429")
	if err != nil {
		panic(err)
	}
}

func SendCharacterMessage(s *discordgo.Session, m *discordgo.MessageCreate, cleanedMessage string, name string, avatar string) {
	if webhook.ChannelID != m.ChannelID {
		s.WebhookEdit(webhook.ID, webhook.Name, webhook.Avatar, m.ChannelID)
	}

	_, err := s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content:     cleanedMessage,
		Username:    name,
		AvatarURL:   avatar,
		Attachments: m.Attachments,
	})
	HandleErr(s, m.ChannelID, err)
}

func SendCharacterReply(s *discordgo.Session, m *discordgo.MessageCreate, cleanedMessage string, name string, avatar string) {
	if webhook.ChannelID != m.ChannelID {
		s.WebhookEdit(webhook.ID, webhook.Name, webhook.Avatar, m.ChannelID)
	}

	var refName string
	msg, err := s.ChannelMessage(m.Reference().ChannelID, m.ReferencedMessage.ID)
	member, err := s.GuildMember(secrets.GuildID, msg.Author.ID)
	if err != nil {
		refName = msg.Author.Username
	} else {
		refName = member.Nick
	}

	_, err = s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Content: fmt.Sprintf("> [Replying to](https://discord.com/channels/@me/%v/%v): %v\n\n%v",
			msg.ChannelID,
			msg.ID,
			refName,
			cleanedMessage,
		),
		Username:    name,
		AvatarURL:   avatar,
		Attachments: m.Attachments,
	})
	HandleErr(s, m.ChannelID, err)
}

func SendPersonalityMessage(s *discordgo.Session, channelID string, msg string, name string, pfpLink string, req *openai.ChatCompletionRequest) {
	if webhook.ChannelID != channelID {
		s.WebhookEdit(webhook.ID, webhook.Name, webhook.Avatar, channelID)
	}

	_, err := s.WebhookExecute(webhook.ID, webhook.Token, false, &discordgo.WebhookParams{
		Username:  name,
		Content:   msg,
		AvatarURL: pfpLink,
	})
	HandleErr(s, channelID, err)
}

func SendPersonalityReply(s *discordgo.Session, m *discordgo.MessageCreate, msg string, name string, pfpLink string, req *openai.ChatCompletionRequest) {
	if webhook.ChannelID != m.ChannelID {
		s.WebhookEdit(webhook.ID, webhook.Name, webhook.Avatar, m.ChannelID)
	}

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
