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
	"BetterScorch/stocks"
	"BetterScorch/webhooks"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var running = false

func main() {
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewPrettyHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

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

	session.AddHandler(threadCreateHandler)
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

	go func() {
		sig := <-sigChan
		slog.Info("Received signal", "signal", sig)

		// Perform any cleanup here if needed
		fmt.Print("|   Reviving all executed members... ")
		execution.ReviveAll(session, "1196943729387372634")
		fmt.Println("Done")
		fmt.Print("|   Cancelling polls... ")
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
		s.GuildMemberRoleAdd("1195135473006420048", "384422339393355786", "1195858179590987866")
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

		go stocks.RegularHandler()
		s.GuildMemberRoleAdd("1195135473006420048", "384422339393355786", "1251675947787096115")

		fmt.Println("Start successful, beginning log")
		fmt.Println("---------------------------------------------------")
		running = true
	}
}

func threadCreateHandler(s *discordgo.Session, tc *discordgo.ThreadCreate) {
	if execution.IsDead(tc.OwnerID) {
		slog.Info("Deleted thread from executed user")
		s.ChannelDelete(tc.ID)
	}
}
