package secretmessages

import (
	"BetterScorch/database"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func SendSecretMessage(s *discordgo.Session, msg string, channelID string, id string, nick string) {
	lastInID, _ := database.Insert("EncryptedMessage", &database.DBValue{Name: "msg", Value: msg})

	embed := &discordgo.MessageEmbed{
		Title: "This is a super secret SWAG message",
		Author: &discordgo.MessageEmbedAuthor{
			Name:    nick,
			URL:     fmt.Sprintf("https://aha-rp.org/get/pilots/%v", strings.ReplaceAll(nick, " ", "")),
			IconURL: fmt.Sprintf("https://aha-rp.org/static/assets/avatars/%v.png", id),
		},
		Color: 16738740,
	}

	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "decrypt" + strconv.Itoa(int(lastInID)),
						Label:    "Decrypt",
					},
				},
			},
		},
	})
}
