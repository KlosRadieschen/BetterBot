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

	if m.Author.Bot || execution.CheckAndDeleteExecuteeMessage(s, m) || Sleeping || (channel.ParentID != "1234128503968891032" && channel.ParentID != "1300423257262133280") {
		if m.Author.Bot {
			execution.CheckAndDeleteExecuteeTupperMessage(s, m, Msgs)
		}
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
	handleMCRef(s, m)
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

func handleMCRef(s *discordgo.Session, m *discordgo.MessageCreate) {
	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("chicken"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("jockey"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# CHICKEN JOCKEY", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}

	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("nether"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# THE NETHER", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}

	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("water"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("bucket"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# WATER BUCKET, RELEASE", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}

	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("flint"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("steel"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# FLINT AND STEEL", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}

	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("craft"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("crafting"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("table"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# THIS, IS A CRAFTING TABLE", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}

	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("mine"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("minecraft"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# FIRST WE MINE THEN WE CRAFT! LETS MINECRAFT", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}

	if regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("villager"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("villagers"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil || regexp.MustCompile(fmt.Sprintf(`\b%s\b`, regexp.QuoteMeta("village"))).FindStringSubmatch(strings.ToLower(m.Content)) != nil {
		sender.SendCharacterMessage(s, m, "# THESE GUYS? THEY’RE THE VILLAGERS!", "Jack \"Steve\" Black", "https://platform.polygon.com/wp-content/uploads/sites/2/2025/04/MCDMIMO_WB046.jpg?quality=90&strip=all&crop=0,0,100,100&w=2400")
	}
}
