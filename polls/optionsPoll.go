package polls

import (
	"BetterScorch/secrets"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type optionPoll struct {
	votes       []int
	voters      map[string][]int
	options     []string
	multioption bool
	cancel      func()
}

const pollChannelID = "1203821534175825942"

var optionPolls = make(map[string]*optionPoll)
var PollMutex sync.Mutex

func CreateOptionsPoll(s *discordgo.Session, creatorID string, multioption bool, endTime time.Time, question string, options ...string) (string, context.Context) {
	emojis := []string{"🔥", "🍷", "💀", "👻", "🎶"}
	votes := []int{}
	ctx, cancelFunc := context.WithTimeout(context.Background(), endTime.Sub(time.Now()))
	optionsString := "\n"

	row := discordgo.ActionsRow{}
	for i, option := range options {
		if len(option) < 75 {
			row.Components = append(row.Components, discordgo.Button{
				CustomID: fmt.Sprintf("pollvote%v", i),
				Label:    option,
				Style:    discordgo.PrimaryButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: emojis[i],
				},
				Disabled: false,
			})
		} else {
			row.Components = append(row.Components, discordgo.Button{
				CustomID: fmt.Sprintf("pollvote%v", i),
				Style:    discordgo.PrimaryButton,
				Emoji: &discordgo.ComponentEmoji{
					Name: emojis[i],
				},
				Disabled: false,
			})
		}

		if len(option) > 15 {
			optionsString += fmt.Sprintf("\n- %v: %v", emojis[i], option)
		}

		votes = append(votes, 0)
	}

	member, _ := s.GuildMember(secrets.GuildID, creatorID)
	var pollTypeString string
	if multioption {
		pollTypeString = "Multi-option poll"
	} else {
		pollTypeString = "Single-option poll"
	}

	pollMsg, err := s.ChannelMessageSendComplex(pollChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("# %v\n(%v by %v)\nPoll expires <t:%v:R>%v",
			question,
			pollTypeString,
			member.Mention(),
			endTime.Unix(),
			optionsString,
		),
		Components: []discordgo.MessageComponent{
			row,
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						CustomID: "pollshow",
						Label:    "Show votes",
						Style:    discordgo.SuccessButton,
						Disabled: false,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	optionPolls[pollMsg.ID] = &optionPoll{votes: votes, voters: make(map[string][]int), options: options, multioption: multioption, cancel: cancelFunc}

	return pollMsg.ID, ctx
}

func WaitAndEvaluate(s *discordgo.Session, pollID string, ctx context.Context) {
	thread, _ := s.MessageThreadStart(pollChannelID, pollID, "Discussion", 60)
	<-ctx.Done()

	slog.Info("Poll ended", "pollType", "Options-poll", "ID", pollID)

	endTime, _ := ctx.Deadline()
	updatePollMessage(s, pollID, true, !time.Now().After(endTime))

	allVotes, _ := GetAllVotesEmbeds(s, pollID)
	s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{Embeds: allVotes})
	s.ChannelEdit(thread.ID, &discordgo.ChannelEdit{Name: "Results"})

	optionPolls[pollID] = nil
}

func AddVote(s *discordgo.Session, id string, voterID string, optionNumber int) error {
	poll, exists := optionPolls[id]
	if !exists {
		return fmt.Errorf("Poll doesn't exist")
	}

	voterNums, hasVoted := poll.voters[voterID]
	if !poll.multioption {
		handleSingleOptionVote(poll, voterID, voterNums, hasVoted, optionNumber)
	} else {
		handleMultiOptionVote(poll, voterID, voterNums, hasVoted, optionNumber)
	}

	updatePollMessage(s, id, false, false)
	return nil
}

func handleSingleOptionVote(poll *optionPoll, voterID string, voterNums []int, hasVoted bool, optionNumber int) {
	if hasVoted {
		if voterNums[0] == optionNumber {
			poll.votes[voterNums[0]]--
			delete(poll.voters, voterID)
		} else {
			poll.votes[voterNums[0]]--
			poll.votes[optionNumber]++
			poll.voters[voterID] = []int{optionNumber}
		}
	} else {
		poll.votes[optionNumber]++
		poll.voters[voterID] = []int{optionNumber}
	}
}

func handleMultiOptionVote(poll *optionPoll, voterID string, voterNums []int, hasVoted bool, optionNumber int) {
	if hasVoted && slices.Contains(voterNums, optionNumber) {
		poll.votes[optionNumber]--
		voterNums = removeElement(voterNums, optionNumber)
	} else {
		poll.votes[optionNumber]++
		voterNums = append(voterNums, optionNumber)
	}
	poll.voters[voterID] = voterNums
}

func removeElement(slice []int, element int) []int {
	index := slices.Index(slice, element)
	if index == -1 {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

func updatePollMessage(s *discordgo.Session, pollID string, isOver bool, isCancelled bool) {
	poll, _ := s.ChannelMessage(pollChannelID, pollID)
	votesSum := GetVotesSum(pollID)

	row := discordgo.ActionsRow{}
	for i, option := range poll.Components[0].(*discordgo.ActionsRow).Components {
		button := option.(*discordgo.Button)
		label := fmt.Sprintf("%v: %v", strings.Split(button.Label, ":")[0], optionPolls[pollID].votes[i])

		if isOver {
			var percentage float32
			if votesSum == 0 {
				percentage = 0
			} else {
				percentage = float32(optionPolls[pollID].votes[i]) / float32(votesSum) * 100
			}
			label = fmt.Sprintf("%v (%.0f%%)", label, percentage)
		}

		row.Components = append(row.Components, discordgo.Button{
			CustomID: button.CustomID,
			Style:    button.Style,
			Label:    label,
			Disabled: isOver,
			Emoji:    button.Emoji,
		})
	}

	var edit discordgo.MessageEdit
	if isOver {
		var newContent string
		if !isCancelled {
			newContent = strings.Replace(poll.Content, "expires", "expired", -1)
		} else {
			newContent = strings.Replace(poll.Content, "expires", "expired", -1) + " (poll ended early)"
		}

		edit = discordgo.MessageEdit{
			Channel: poll.ChannelID,
			ID:      poll.ID,
			Content: &newContent,
			Embeds:  &poll.Embeds,
			Components: &[]discordgo.MessageComponent{
				row,
			},
		}
	} else {
		edit = discordgo.MessageEdit{
			Channel: poll.ChannelID,
			ID:      poll.ID,
			Content: &poll.Content,
			Embeds:  &poll.Embeds,
			Components: &[]discordgo.MessageComponent{
				row,
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							CustomID: "pollshow",
							Label:    "Show votes",
							Style:    discordgo.SuccessButton,
							Disabled: false,
						},
					},
				},
			},
		}
	}

	s.ChannelMessageEditComplex(&edit)
}

func GetAllVotesString(pollID string) (string, error) {
	_, exists := optionPolls[pollID]
	if !exists {
		return "", fmt.Errorf("Poll doesn't exist")
	}

	emojis := []string{"🔥", "🍷", "💀", "👻", "🎶"}
	allVotes := ""
	for i := range len(optionPolls[pollID].votes) {
		voterPings := []string{}
		for name, votes := range optionPolls[pollID].voters {
			if slices.Contains(votes, i) {
				voterPings = append(voterPings, fmt.Sprintf("<@%v>", name))
			}
		}

		if len(voterPings) != 0 {
			allVotes += fmt.Sprintf("%v: %v\n\n", emojis[i], strings.Join(voterPings, ", "))
		}
	}

	if allVotes == "" {
		return "There are no votes", nil
	}
	return allVotes, nil
}

func GetAllVotesEmbeds(s *discordgo.Session, pollID string) ([]*discordgo.MessageEmbed, error) {
	emojis := []string{"🔥", "🍷", "💀", "👻", "🎶"}
	poll, exists := optionPolls[pollID]
	if !exists {
		return nil, fmt.Errorf("Poll doesn't exist")
	}

	embeds := []*discordgo.MessageEmbed{}

	for i := range len(poll.votes) {
		embeds = append(embeds, &discordgo.MessageEmbed{Title: emojis[i] + ": " + poll.options[i], Color: 0x3498db})
	}

	for i := range poll.votes {
		for voter, votes := range poll.voters {
			if slices.Contains(votes, i) {
				member, _ := s.GuildMember(secrets.GuildID, voter)
				embeds[i].Fields = append(embeds[i].Fields, &discordgo.MessageEmbedField{Name: member.Nick})
			}
		}
	}

	var filteredEmbeds []*discordgo.MessageEmbed
	for _, embed := range embeds {
		if embed.Fields != nil {
			filteredEmbeds = append(filteredEmbeds, embed)
		}
	}

	return filteredEmbeds, nil
}

func GetVotesSum(pollID string) int {
	sum := 0
	for _, vote := range optionPolls[pollID].votes {
		sum += vote
	}
	return sum
}
