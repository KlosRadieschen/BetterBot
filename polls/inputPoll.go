package polls

import (
	"BetterScorch/secrets"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type inputPoll struct {
	votes       map[string][]string
	multioption bool
}

var inputPolls = make(map[string]*inputPoll)

func CreateInputPoll(s *discordgo.Session, creatorID string, multioption bool, endTime time.Time, question string) string {
	member, _ := s.GuildMember(secrets.GuildID, creatorID)

	var pollTypeString string
	if multioption {
		pollTypeString = "Multi-input poll"
	} else {
		pollTypeString = "Single-input poll"
	}

	pollMsg, _ := s.ChannelMessageSendComplex(pollChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("# %v\n(%v by %v)\nPoll expires <t:%v:R>",
			question,
			pollTypeString,
			member.Mention(),
			endTime.Unix(),
		),
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "inputPollModalCreate",
						Label:    "Respond",
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "inputpollshow",
						Label:    "Show responses",
						Style:    discordgo.SuccessButton,
					},
				},
			},
		},
	})

	inputPolls[pollMsg.ID] = &inputPoll{votes: make(map[string][]string), multioption: multioption}

	return pollMsg.ID
}

func SubmitInputPollResponse(s *discordgo.Session, id string, voterID string, response string) error {
	poll, exists := inputPolls[id]
	if !exists {
		return fmt.Errorf("Poll doesn't exist")
	}

	if poll.multioption {
		poll.votes[voterID] = append(poll.votes[voterID], response)
	} else {
		poll.votes[voterID] = []string{response}
	}
	return nil
}

func WaitAndEvaluateInput(s *discordgo.Session, pollID string, endTime time.Time) {
	time.Sleep(endTime.Sub(time.Now()))

	poll, _ := s.ChannelMessage(pollChannelID, pollID)
	newContent := strings.Replace(poll.Content, "expires", "expired", -1)
	edit := discordgo.MessageEdit{
		Channel:    poll.ChannelID,
		ID:         poll.ID,
		Content:    &newContent,
		Components: &[]discordgo.MessageComponent{},
	}

	s.ChannelMessageEditComplex(&edit)

	thread, _ := s.MessageThreadStart(pollChannelID, pollID, "Results", 60)
	embeds, _ := GetAllInputsEmbeds(s, pollID)
	s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{Embeds: embeds})
	trueBool := true
	s.ChannelEdit(thread.ID, &discordgo.ChannelEdit{Locked: &trueBool})

	inputPolls[pollID] = nil
}

func GetAllInputsEmbeds(s *discordgo.Session, pollID string) ([]*discordgo.MessageEmbed, error) {
	_, exists := inputPolls[pollID]
	if !exists {
		return nil, fmt.Errorf("Poll doesn't exist")
	}

	embeds := []*discordgo.MessageEmbed{}

	embedNum := int(math.Ceil(float64(GetNumberOfInputs(pollID)) / 25))
	for range embedNum {
		embeds = append(embeds, &discordgo.MessageEmbed{Color: 0x3498db})
	}

	count := 0
	for voter, responses := range inputPolls[pollID].votes {
		for _, response := range responses {
			member, _ := s.GuildMember(secrets.GuildID, voter)
			embeds[int(math.Ceil(float64(count+1)/25))-1].Fields = append(embeds[int(math.Ceil(float64(count+1)/25))-1].Fields, &discordgo.MessageEmbedField{
				Name:   member.Nick + ":",
				Value:  response,
				Inline: false,
			})
			count++
		}
	}

	return embeds, nil
}

func GetNumberOfInputs(pollID string) int {
	sum := 0
	for _, votes := range inputPolls[pollID].votes {
		for range votes {
			sum++
		}
	}
	return sum
}
