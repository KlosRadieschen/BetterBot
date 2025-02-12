package execution

import (
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"fmt"
	"slices"

	"github.com/bwmarrin/discordgo"
)

type executee struct {
	id         string
	count      int
	role       string
	sacrificed bool
}

var executees = []executee{}
var roles = []string{"1195135956471255140", "1195858311627669524", "1195858271349784639", "1195136106811887718", "1195858179590987866", "1195137362259349504", "1195136284478410926", "1195137253408768040", "1250582641921757335", "1195758308519325716", "1195758241221722232", "1195758137563689070", "1195757362439528549", "1195136491148550246", "1195708423229165578", "1195137477497868458", "1195136604373782658", "1248776818664935525", "1277091839647678575"}

func Execute(s *discordgo.Session, userID string, channelID string, sacrificed bool) {
	for _, executee := range executees {
		if executee.id == userID {
			index := slices.Index(executees, executee) //necessary to edit the executee directly
			executees[index].count++
			executees[index].sacrificed = false

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
	})
	sender.SendMessage(s, channelID, fmt.Sprintf("%v is fucking dead", Member(s, userID).Mention()))
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

	sender.SendMessage(s, channelID, fmt.Sprintf("%v has been revived!", member.Mention()))
}

func IsDead(userID string) bool {
	for _, executee := range executees {
		if executee.id == userID {
			return true
		}
	}
	return false
}

func IsSacrificed(userID string) bool {
	for _, executee := range executees {
		if executee.id == userID {
			return executee.sacrificed
		}
	}
	return false
}

func CheckAndDeleteExecuteeMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if IsDead(m.Author.ID) {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		return true
	}
	return false
}

func Member(s *discordgo.Session, userID string) *discordgo.Member {
	member, _ := s.GuildMember(secrets.GuildID, userID)
	return member
}
