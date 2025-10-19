package messages

import (
	"BetterScorch/ai"
	"BetterScorch/execution"
	"BetterScorch/sender"
	"BetterScorch/webhooks"
	"container/list"
	"fmt"
	"log"
	"log/slog"
	"regexp"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type UserMessage struct {
	SenderName  string
	SenderPFP   string
	RecipientID string
	Message     string
}

var Sleeping = false
var Msgs = make(map[string]*list.List)
var UserMessages = []UserMessage{}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	defer func() {
		if r := recover(); r != nil {
			s.ChannelMessageSendComplex("1196943729387372634", &discordgo.MessageSend{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title: fmt.Sprintf("PANIC: %v", r),
						Color: 0xFF69B4,
					},
				},
			})
			slog.Error(fmt.Sprintf("%v", r))
			fmt.Println(string(debug.Stack()))
		}
	}()

	channel, _ := s.Channel(m.ChannelID)

	if Msgs[m.Author.ID] == nil {
		Msgs[m.Author.ID] = list.New()
	}
	Msgs[m.Author.ID].PushBack(m.Message)
	if Msgs[m.Author.ID].Len() > 5 {
		Msgs[m.Author.ID].Remove(Msgs[m.Author.ID].Front())
	}
	go webhooks.CheckAndUseCharacters(s, m)
	go checkAndSendUserMessages(s, m)

	if m.Author.Bot {
		execution.CheckAndDeleteExecuteeTupperMessage(s, m, Msgs)
	}
	if m.Author.Bot || execution.CheckAndDeleteExecuteeMessage(s, m) || Sleeping || (channel.ParentID != "1234128503968891032" && channel.ParentID != "1300423257262133280") {
		return
	}

	go webhooks.CheckAndRespondPersonalities(s, m)
	go handleAIResponses(s, m)

	for _, response := range responses {
		for _, trigger := range response.triggers {
			// Magical RegEx bullshiterry
			match := regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta(trigger))).FindStringSubmatch(strings.ToLower(m.Content))
			if match != nil {
				log.Println("Found trigger: " + trigger)
				response.handleResponse(s, m)
			}
		}
	}
}

func handleAIResponses(s *discordgo.Session, m *discordgo.MessageCreate) {
	scorchUser, _ := s.User("1196526025211904110")
	found := slices.ContainsFunc(m.Mentions, func(u *discordgo.User) bool {
		return u.ID == scorchUser.ID
	})

	if (m.Type == 19 && m.ReferencedMessage.Author.ID == "1196526025211904110") || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("scorch"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || found {
		log.Println("AI response triggered: " + m.Content)

		s.ChannelTyping(m.ChannelID)
		resp, embed, err := ai.GenerateResponse(m.Member.Nick, m.Content)

		if !sender.HandleErr(s, m.ChannelID, err) {
			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Reference: m.Reference(),
				Content:   resp,
				Embed:     embed,
			})
		}

		if embed != nil && embed.Title == "Used /execute" {
			execution.Execute(s, m.Author.ID, m.ChannelID, false)
		}
	}
}

func checkOOCChannel(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		return false
	}
	return channel.ParentID == "1234128503968891032" || channel.ParentID == "1300423257262133280"
}

func checkAndSendUserMessages(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, userMessage := range UserMessages {
		if m.Author.ID == userMessage.RecipientID {
			embed := discordgo.MessageEmbed{
				Description: userMessage.Message,
				Color:       0xFF69B4,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    userMessage.SenderName,
					IconURL: userMessage.SenderPFP,
				},
			}

			s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Reference: m.Reference(), Embeds: []*discordgo.MessageEmbed{&embed}})
			UserMessages = slices.DeleteFunc(UserMessages, func(um UserMessage) bool {
				return um.SenderName == userMessage.SenderName
			})
		}
	}
}
