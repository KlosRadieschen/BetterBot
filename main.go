package main

import (
	"BetterScorch/ai"
	"BetterScorch/commands"
	"BetterScorch/database"
	"BetterScorch/execution"
	"BetterScorch/messages"
	"BetterScorch/polls"
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"BetterScorch/webhooks"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var running = false

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

	// Create channels for signals
	sigChan := make(chan os.Signal, 1)

	// Register the signals we want to handle
	signal.Notify(sigChan,
		os.Interrupt, // Keyboard interrupt (Ctrl+C)
		syscall.SIGTERM,
		syscall.SIGQUIT) // Systemd service stop

	// Run in a goroutine so it doesn't block the main program
	go func() {
		sig := <-sigChan
		fmt.Printf("\nReceived signal: %v\n", sig)

		// Perform any cleanup here if needed
		fmt.Print("|   Reviving all executed members...")
		execution.ReviveAll(session, "1196943729387372634")
		fmt.Println("Done")
		fmt.Print("|   Cancelling polls...")
		polls.CancelAll()
		fmt.Println("Done")
		session.ChannelMessageSendComplex("1196943729387372634", &discordgo.MessageSend{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title: "Better Scorch shutting down",
					Color: 0xFF69B4,
				},
			},
		})
		fmt.Println("Shutdown complete")

		os.Exit(0) // Exit with status code 0
	}()

	<-make(chan struct{})
}

func readyHandler(s *discordgo.Session, r *discordgo.Ready) {
	if !running {
		fmt.Println("|   Initialising webhooks")
		fmt.Print("    |   Loading Scorch webhook... ")
		sender.InitWebhook(s)
		fmt.Println("Done")
		fmt.Print("    |   Loading characters... ")
		webhooks.RetrieveCharacters()
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
		running = true
	}
}
