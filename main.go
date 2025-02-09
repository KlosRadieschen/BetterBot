package main

import (
	"BetterScorch/ai"
	"BetterScorch/commands"
	"BetterScorch/messages"
	"BetterScorch/secrets"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Commencing startup sequence")
	fmt.Println("|   Initialising AI package")
	ai.Init()

	fmt.Println("Initialising session")
	var err error
	session, _ := discordgo.New("Bot " + secrets.BotToken)
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("|   Initialising commands package")
		commands.AddAllCommands(session)
		session.AddHandler(messages.HandleMessage)
		fmt.Println("Start successful, beginning log")
		fmt.Println("---------------------------------------------------")
	})
	err = session.Open()
	if err != nil {
		panic("Couldnt open session")
	}

	session.UpdateListeningStatus("the screams of burning PHC pilots")

	<-make(chan struct{})
}
