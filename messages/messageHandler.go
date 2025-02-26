package messages

import (
	"BetterScorch/ai"
	"BetterScorch/execution"
	"BetterScorch/sender"
	"BetterScorch/webhooks"
	"container/list"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Sleeping = false
var Msgs = make(map[string]*list.List)

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, _ := s.Channel(m.ChannelID)

	if Msgs[m.Author.ID] == nil {
		Msgs[m.Author.ID] = list.New()
	}
	Msgs[m.Author.ID].PushBack(m.Message)
	if Msgs[m.Author.ID].Len() > 5 {
		Msgs[m.Author.ID].Remove(Msgs[m.Author.ID].Front())
	}
	go webhooks.CheckAndUseCharacters(s, m)

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
	if m.Type == 19 && m.ReferencedMessage.Author.ID == "1196526025211904110" || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("scorch"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
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
