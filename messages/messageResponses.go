package messages

import (
	"BetterScorch/sender"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

type messageResponse struct {
	triggers []string
	response string
	isMedia  bool
}

func (mr *messageResponse) handleResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !mr.isMedia {
		s.ChannelMessageSendReply(m.ChannelID, mr.response, m.Reference())
	} else {
		file, err := os.Open(fmt.Sprintf("media/%v.png", mr.response))
		extension := ".png"
		if err != nil {
			file, err = os.Open(fmt.Sprintf("media/%v.mp4", mr.response))
			extension = ".mp4"
			sender.HandleErr(s, m.ChannelID, err)
		}
		defer file.Close()
		reader := discordgo.File{
			Name:   mr.response + extension,
			Reader: file,
		}
		messageContent := &discordgo.MessageSend{
			Files:     []*discordgo.File{&reader},
			Reference: m.Reference(),
		}
		s.ChannelMessageSendComplex(m.ChannelID, messageContent)
	}

}
