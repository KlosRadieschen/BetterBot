package polls

import (
	"BetterScorch/secrets"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type poll struct {
	votes       []int
	voters      map[string][]int
	multioption bool
}

const pollChannelID = "1196943729387372634"

var polls = make(map[string]*poll)

func CreatePoll(s *discordgo.Session, creatorID string, multioption bool, endTime time.Time, question string, options ...string) string {
	emojis := []string{"🔥", "🍷", "💀", "👻", "🎶"}
	votes := []int{}

	components := []discordgo.MessageComponent{}

	for i, option := range options {
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: fmt.Sprintf("pollvote%v", i),
					Label:    option,
					Style:    discordgo.PrimaryButton,
					Emoji: &discordgo.ComponentEmoji{
						Name: emojis[i],
					},
					Disabled: false,
				},
			},
		})
		votes = append(votes, 0)
	}

	member, _ := s.GuildMember(secrets.GuildID, creatorID)
	pollMsg, _ := s.ChannelMessageSendComplex(pollChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("**%v** (by %v)\nPoll expires <t:%v:R>",
			question,
			member.Mention(),
			endTime.Unix(),
		),
		Components: components,
	})

	polls[pollMsg.ID] = &poll{votes: votes, voters: make(map[string][]int), multioption: multioption}

	return pollMsg.ID
}

func WaitAndEvaluate(s *discordgo.Session, pollID string, endTime time.Time) {
	time.Sleep(endTime.Sub(time.Now()))
	updatePollMessage(s, pollID, true)
	polls[pollID] = nil
}

func AddVote(s *discordgo.Session, id string, voterID string, optionNumber int) error {
	poll, exists := polls[id]
	if !exists {
		return fmt.Errorf("Poll doesn't exist")
	}

	voterNums, hasVoted := poll.voters[voterID]
	if !poll.multioption {
		handleSingleOptionVote(poll, voterID, voterNums, hasVoted, optionNumber)
	} else {
		handleMultiOptionVote(poll, voterID, voterNums, hasVoted, optionNumber)
	}

	updatePollMessage(s, id, false)
	return nil
}

func handleSingleOptionVote(poll *poll, voterID string, voterNums []int, hasVoted bool, optionNumber int) {
	if hasVoted && len(voterNums) > 0 {
		poll.votes[voterNums[0]]--
	}
	poll.votes[optionNumber]++
	poll.voters[voterID] = []int{optionNumber}
}

func handleMultiOptionVote(poll *poll, voterID string, voterNums []int, hasVoted bool, optionNumber int) {
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

func updatePollMessage(s *discordgo.Session, id string, isOver bool) {
	poll, _ := s.ChannelMessage(pollChannelID, id)
	votesSum := getVotesSum(id)

	var components []discordgo.MessageComponent
	for i, option := range poll.Components {
		button := option.(*discordgo.ActionsRow).Components[0].(*discordgo.Button)
		label := fmt.Sprintf("%v: %v", strings.Split(button.Label, ":")[0], polls[id].votes[i])

		var percentage int
		if isOver {
			if votesSum == 0 {
				percentage = 0
			} else {
				percentage = polls[id].votes[i] / votesSum * 100
			}
			label = fmt.Sprintf("%v (%v)", label, percentage)
		}

		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: button.CustomID,
					Style:    button.Style,
					Label:    fmt.Sprintf("%v: %v", strings.Split(button.Label, ":")[0], polls[id].votes[i]),
					Disabled: false,
					Emoji:    button.Emoji,
				},
			},
		})
	}

	newContent := poll.Content
	if isOver {
		newContent = strings.Replace(poll.Content, "expires", "expired", -1)
	}

	edit := discordgo.MessageEdit{
		Channel:    poll.ChannelID,
		ID:         poll.ID,
		Content:    &newContent,
		Embeds:     &poll.Embeds,
		Components: &components,
	}

	s.ChannelMessageEditComplex(&edit)
}

func getVotesSum(pollID string) int {
	sum := 0
	for _, vote := range polls[pollID].votes {
		sum += vote
	}
	return sum
}
