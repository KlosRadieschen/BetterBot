package messages

import (
	"github.com/bwmarrin/discordgo"
)

type messageResponseText struct {
	trigger  []string
	response string
}

type messageResponseMedia struct {
	trigger  []string
	response discordgo.File
}

type messageResponse interface {
	triggers() []string
	handle(s *discordgo.Session, m *discordgo.MessageCreate)
}

func (mrt messageResponseText) triggers() []string {
	return mrt.trigger
}

func (mrt *messageResponseText) handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSendReply(m.ChannelID, mrt.response, m.Reference())
}

func (mrm messageResponseMedia) triggers() []string {
	return mrm.trigger
}

func (mrm *messageResponseMedia) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Implement

}

func getTriggers(resp messageResponse) []string {
	return resp.triggers()
}

func handleResponse(s *discordgo.Session, m *discordgo.MessageCreate, resp messageResponse) {
	resp.handle(s, m)
}
