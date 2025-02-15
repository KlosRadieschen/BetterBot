package main

import (
	"BetterScorch/ai"
	"BetterScorch/commands"
	"BetterScorch/database"
	"BetterScorch/messages"
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("Commencing startup sequence")

	fmt.Print("|   Initialising AI package... ")
	ai.Init()
	fmt.Println("Done")

	fmt.Print("|   Opening database connection... ")
	database.Connect()
	fmt.Println("Done")

	fmt.Print("|   Initialising session... ")
	var err error
	session, _ := discordgo.New("Bot " + secrets.BotToken)
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	session.AddHandler(readyHandler)
	err = session.Open()
	if err != nil {
		panic("Couldnt open session")
	}
	fmt.Println("Done")

	session.UpdateListeningStatus("the screams of burning PHC pilots")

	<-make(chan struct{})
}

func readyHandler(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Print("|   Initialising webhook... ")
	sender.InitWebhook(s)
	fmt.Println("Done")
	fmt.Println("|   Initialising commands package")
	commands.AddAllCommands(s)
	s.AddHandler(messages.HandleMessage)

	s.ChannelMessageSendComplex("1196943729387372634", &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Better Scorch started",
				Color: 0xFF69B4,
			},
		},
	})

	fmt.Println("Start successful, beginning log")
	fmt.Println("---------------------------------------------------")
}
