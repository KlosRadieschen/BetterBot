package commands

import (
	"BetterScorch/secrets"
	"fmt"
	"log"
	"slices"

	"github.com/bwmarrin/discordgo"
)

type command struct {
	declaration *discordgo.ApplicationCommand
	handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func AddAllCommands(s *discordgo.Session) {
	fmt.Print("    |   Deleting current commands... ")
	// s.ApplicationCommandBulkOverwrite(s.State.Application.ID, secrets.GuildID, []*discordgo.ApplicationCommand{})
	fmt.Println("Done")

	fmt.Print("    |   Re-adding existing commands... ")
	for i, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.Application.ID, secrets.GuildID, command.declaration)
		if err != nil {
			panic(err.Error())
		}

		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.ApplicationCommandData().Name == command.declaration.Name {
				log.Println("Received Command: " + i.ApplicationCommandData().Name)
				command.handler(s, i)
			}
		})
		fmt.Println()
		fmt.Printf("        |   %.2f/100", (float32(i) / float32(len(commands)) * 100))
	}
	fmt.Println()
	fmt.Println("Done")
}

func isHC(m *discordgo.Member) bool {
	if IsKlos(m) {
		return true
	}

	var roles = []string{"1195135956471255140", "1195858311627669524", "1195858271349784639", "1195136106811887718"}

	for _, role := range m.Roles {
		if slices.Contains(roles, role) {
			return true
		}
	}
	return false
}

func IsKlos(m *discordgo.Member) bool {
	return m.User.ID == "384422339393355786"
}
