package main

import (
	"BetterScorch/ai"
	"BetterScorch/messages"
	"BetterScorch/secrets"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	ai.Init()
	messages.Init()

	var err error
	session, _ := discordgo.New("Bot " + secrets.BotToken)

	session.AddHandler(messages.HandleMessage)

	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
		fmt.Println()
	})
	err = session.Open()
	if err != nil {
		panic("Couldnt open session")
	}

	session.UpdateListeningStatus("the screams of burning PHC pilots")

	<-make(chan struct{})
}
