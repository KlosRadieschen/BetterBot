package execution

import (
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"BetterScorch/stocks"
	"container/list"
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type executee struct {
	id         string
	count      int
	role       string
	sacrificed bool
	startTime  time.Time
}

var executees = []executee{}
var roles = []string{"1195135956471255140", "1195858311627669524", "1195858271349784639", "1195136106811887718", "1195858179590987866", "1195137362259349504", "1195136284478410926", "1195137253408768040", "1250582641921757335", "1195758308519325716", "1195758241221722232", "1195758137563689070", "1195757362439528549", "1195136491148550246", "1195708423229165578", "1195137477497868458", "1195136604373782658", "1248776818664935525", "1277091839647678575"}

func Execute(s *discordgo.Session, userID string, channelID string, sacrificed bool) {
	for _, executee := range executees {
		if executee.id == userID {
			index := slices.Index(executees, executee) //necessary to edit the executee directly
			executees[index].count++
			executees[index].sacrificed = false

			stocks.ModifyCompanyValue("Execution Solutions LLC", -500)
			sender.SendMessage(s, channelID, fmt.Sprintf("Increasing %v's execution count to %v!", Member(s, userID).Mention(), executee.count+1))
			return
		}
	}

	roleID := ""
	for _, role := range Member(s, userID).Roles {
		if slices.Contains(roles, role) {
			roleID = role
			s.GuildMemberRoleRemove(secrets.GuildID, userID, roleID)
		}
	}
	s.GuildMemberRoleAdd(secrets.GuildID, userID, "1253410294999548046")

	executees = append(executees, executee{
		id:         Member(s, userID).User.ID,
		count:      1,
		role:       roleID,
		sacrificed: sacrificed,
		startTime:  time.Now(),
	})
	stocks.ModifyCompanyValue("Execution Solutions LLC", -500)
	sender.SendMessage(s, channelID, fmt.Sprintf("%v is fucking dead", Member(s, userID).Mention()))
}

func GambleExecute(s *discordgo.Session, i *discordgo.InteractionCreate, attackerID string, victimID string) {
	attackerMember, _ := s.GuildMember(i.GuildID, attackerID)
	victimMember, _ := s.GuildMember(i.GuildID, victimID)

	var msg *discordgo.Message
	if victimID == "942159289836011591" {
		msg, _ = s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(attackerMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(victimMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(victimMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(victimMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: "BOTH FUCKING DIE",
				},
				{
					Color: 0xFF69B4,
					Title: "NOBODY FUCKING DIES",
				},
			},
		})
	} else if !IsAdminAbuser(victimMember) {
		msg, _ = s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(attackerMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(victimMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: "BOTH FUCKING DIE",
				},
				{
					Color: 0xFF69B4,
					Title: "NOBODY FUCKING DIES",
				},
			},
		})
	} else {
		msg, _ = s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(attackerMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(attackerMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(attackerMember.Nick) + " FUCKING DIES",
				},
				{
					Color: 0xFF69B4,
					Title: strings.ToUpper(attackerMember.Nick) + " FUCKING DIES",
				},
			},
		})
	}

	max := len(msg.Embeds) - 1
	for range max {
		counter := 5
		for counter != 0 {
			countdownString := fmt.Sprintf("ELIMINATING ONE OPTION IN %v", counter)
			msg, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      msg.ID,
				Channel: msg.ChannelID,
				Content: &countdownString,
			})
			counter--
			time.Sleep(1 * time.Second)
		}

		rng := rand.Intn(max + 1)
		max--
		countdownString := fmt.Sprintf("ELIMINATING ONE OPTION IN %v", counter)
		embeds := msg.Embeds
		embeds = append(embeds[:rng], embeds[rng+1:]...)
		msg, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      msg.ID,
			Channel: msg.ChannelID,
			Content: &countdownString,
			Embeds:  &embeds,
		})
	}

	emptyString := ""
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      msg.ID,
		Channel: msg.ChannelID,
		Content: &emptyString,
		Embeds:  &msg.Embeds,
	})

	switch strings.Split(msg.Embeds[0].Title, " FUCKING")[0] {
	case strings.ToUpper(attackerMember.Nick):
		stocks.ModifyCompanyValue("Execution Solutions LLC", 1000)
		Execute(s, attackerID, msg.ChannelID, false)
	case strings.ToUpper(victimMember.Nick):
		stocks.ModifyCompanyValue("Execution Solutions LLC", 1000)
		Execute(s, victimID, msg.ChannelID, false)
	case "BOTH":
		stocks.ModifyCompanyValue("Execution Solutions LLC", 2000)
		Execute(s, attackerID, msg.ChannelID, false)
		Execute(s, victimID, msg.ChannelID, false)
	}
}

func Revive(s *discordgo.Session, userID string, channelID string) {
	member := Member(s, userID)
	for i, executee := range executees {
		if executee.id == member.User.ID {
			for range executee.count {
				sender.SendMessage(s, channelID, fmt.Sprintf("%v\nhttps://tenor.com/view/cat-revive-friends-animated-friendship-gif-8246087956711984034", member.Mention()))
			}
			s.GuildMemberRoleRemove(secrets.GuildID, userID, "1253410294999548046")
			s.GuildMemberRoleAdd(secrets.GuildID, userID, executee.role)
			executees = append(executees[:i], executees[i+1:]...)
			break
		}
	}

	stocks.ModifyCompanyValue("Revival Technologies", -500)
	sender.SendMessage(s, channelID, fmt.Sprintf("%v has been revived!", member.Mention()))
}

func GambleRevive(s *discordgo.Session, i *discordgo.InteractionCreate, sacrificerID string, victimID string) {
	sacrificerMember, _ := s.GuildMember(i.GuildID, sacrificerID)
	victimMember, _ := s.GuildMember(i.GuildID, victimID)

	msg, _ := s.ChannelMessageSendComplex(i.ChannelID, &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Color: 0xFF69B4,
				Title: strings.ToUpper(sacrificerMember.Nick) + " FUCKING DIES",
			},
			{
				Color: 0xFF69B4,
				Title: strings.ToUpper(victimMember.Nick) + " FUCKING LIVES",
			},
			{
				Color: 0xFF69B4,
				Title: "SOUL FUCKING TRADE",
			},
			{
				Color: 0xFF69B4,
				Title: "NOTHING FUCKING HAPPENS",
			},
		},
	})

	max := len(msg.Embeds) - 1
	for range max {
		counter := 5
		for counter != 0 {
			countdownString := fmt.Sprintf("ELIMINATING ONE OPTION IN %v", counter)
			msg, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      msg.ID,
				Channel: msg.ChannelID,
				Content: &countdownString,
			})
			counter--
			time.Sleep(1 * time.Second)
		}

		rng := rand.Intn(max + 1)
		max--
		countdownString := fmt.Sprintf("ELIMINATING ONE OPTION IN %v", counter)
		embeds := msg.Embeds
		embeds = append(embeds[:rng], embeds[rng+1:]...)
		msg, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
			ID:      msg.ID,
			Channel: msg.ChannelID,
			Content: &countdownString,
			Embeds:  &embeds,
		})
	}

	emptyString := ""
	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:      msg.ID,
		Channel: msg.ChannelID,
		Content: &emptyString,
		Embeds:  &msg.Embeds,
	})

	switch strings.Split(msg.Embeds[0].Title, " FUCKING")[0] {
	case strings.ToUpper(sacrificerMember.Nick):
		stocks.ModifyCompanyValue("Execution Solutions LLC", 1000)
		Execute(s, sacrificerID, msg.ChannelID, false)
	case strings.ToUpper(victimMember.Nick):
		stocks.ModifyCompanyValue("Revival Technologies", 2000)
		Revive(s, victimID, msg.ChannelID)
	case "SOUL":
		stocks.ModifyCompanyValue("Revival Technologies", 2000)
		stocks.ModifyCompanyValue("Execution Solutions LLC", 1000)
		Execute(s, sacrificerID, msg.ChannelID, false)
		Revive(s, victimID, msg.ChannelID)
	}
}

func ReviveAll(s *discordgo.Session, channelID string) {
	for i, executee := range executees {
		member := Member(s, executee.id)
		for range executee.count {
			sender.SendMessage(s, channelID, fmt.Sprintf("%v\nhttps://tenor.com/view/cat-revive-friends-animated-friendship-gif-8246087956711984034", member.Mention()))
		}
		s.GuildMemberRoleRemove(secrets.GuildID, executee.id, "1253410294999548046")
		s.GuildMemberRoleAdd(secrets.GuildID, executee.id, executee.role)
		executees = append(executees[:i], executees[i+1:]...)
	}

	sender.SendMessage(s, channelID, "Everyone has been revived!")
}

func IsDead(userID string) bool {
	return getExecutee(userID) != nil
}

func IsSacrificed(userID string) bool {
	e := getExecutee(userID)
	if e != nil {
		return e.sacrificed
	}
	return false
}

func getExecutee(userID string) *executee {
	for _, executee := range executees {
		if executee.id == userID {
			return &executee
		}
	}
	return nil
}

func CheckAndDeleteExecuteeMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if IsDead(m.Author.ID) {
		if time.Since(getExecutee(m.Author.ID).startTime) > 30*time.Minute {
			slog.Info("Automatic revive triggered", "userID", m.Author.ID)
			sender.SendMessage(s, "1196943729387372634", "Automatically reviving "+Member(s, m.Author.ID).Mention())
			Revive(s, m.Author.ID, "1196943729387372634")
		} else {
			slog.Info("Deleted message from executed user", "userID", m.Author.ID, "message", m.Content)
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			return true
		}
	}
	return false
}

func CheckAndDeleteExecuteeTupperMessage(s *discordgo.Session, m *discordgo.MessageCreate, msgs map[string]*list.List) {
	for _, executee := range executees {
		if msgs[executee.id] != nil && m.Content != "" {
			// Split the current message content into lines
			lines := strings.Split(m.Content, "\n")

			var filteredContentBuilder strings.Builder

			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if !strings.HasPrefix(trimmedLine, ">") {
					if filteredContentBuilder.Len() > 0 {
						filteredContentBuilder.WriteString("\n")
					}
					filteredContentBuilder.WriteString(trimmedLine)
				}
			}

			// Get the filtered content as a single string
			filteredContent := strings.ToLower(filteredContentBuilder.String())

			lastMessage := msgs[executee.id].Back().Value.(*discordgo.Message)

			if strings.Contains(strings.ToLower(lastMessage.Content), filteredContent) {
				slog.Info("Deleted tupper message from executed user", "userID", executee.id, "message", m.Content)
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				return
			}
		}
	}
}

func Member(s *discordgo.Session, userID string) *discordgo.Member {
	member, _ := s.GuildMember(secrets.GuildID, userID)
	return member
}

func IsAdminAbuser(m *discordgo.Member) bool {
	return m.User.ID == "920342100468436993" || m.User.ID == "1079774043684745267" || m.User.ID == "952145898824138792"
}
