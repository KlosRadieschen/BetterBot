package polls

import (
	"BetterScorch/secrets"
	"context"
	"fmt"
	"log/slog"
	"math"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type inputPoll struct {
	votes       map[string][]string
	options     []string
	multioption bool
	cancel      func()
}

var inputPolls = make(map[string]*inputPoll)

func CreateInputPoll(s *discordgo.Session, creatorID string, multioption bool, endTime time.Time, question string) (string, context.Context) {
	member, _ := s.GuildMember(secrets.GuildID, creatorID)
	ctx, cancelFunc := context.WithTimeout(context.Background(), endTime.Sub(time.Now()))

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

	inputPolls[pollMsg.ID] = &inputPoll{votes: make(map[string][]string), multioption: multioption, cancel: cancelFunc}

	return pollMsg.ID, ctx
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

func WaitAndEvaluateInput(s *discordgo.Session, pollID string, ctx context.Context) {
	thread, _ := s.MessageThreadStart(pollChannelID, pollID, "Discussion", 60)

	<-ctx.Done()

	slog.Info("Poll ended", "pollType", "Input-poll", "ID", pollID)

	endTime, _ := ctx.Deadline()
	poll, _ := s.ChannelMessage(pollChannelID, pollID)

	var newContent string
	if time.Now().After(endTime) {
		newContent = strings.Replace(poll.Content, "expires", "expired", -1)
	} else {
		newContent = strings.Replace(poll.Content, "expires", "expired", -1) + " (poll ended early)"
	}

	edit := discordgo.MessageEdit{
		Channel:    poll.ChannelID,
		ID:         poll.ID,
		Content:    &newContent,
		Components: &[]discordgo.MessageComponent{},
	}

	s.ChannelMessageEditComplex(&edit)

	embeds, _ := GetAllInputsEmbeds(s, pollID)
	s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{Embeds: embeds})
	s.ChannelEdit(thread.ID, &discordgo.ChannelEdit{Name: "Results"})

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

func CancelAll() {
	for _, poll := range inputPolls {
		poll.cancel()
		time.Sleep(3 * time.Second)
	}
	for _, poll := range optionPolls {
		poll.cancel()
		time.Sleep(3 * time.Second)
	}
}
