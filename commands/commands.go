package commands

import (
	"BetterScorch/execution"
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"fmt"
	"log"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type command struct {
	declaration *discordgo.ApplicationCommand
	handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func AddAllCommands(s *discordgo.Session) {
	fmt.Print("    |   Adding slash commands... ")

	commandSlice := []*discordgo.ApplicationCommand{}
	for _, command := range commands {
		commandSlice = append(commandSlice, command.declaration)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		defer func() {
			if r := recover(); r != nil {
				sender.Followup(s, i, fmt.Sprintf("PANIC: %v", r))
				sender.Respond(s, i, fmt.Sprintf("PANIC: %v", r), nil)
				fmt.Println(string(debug.Stack()))
			}
		}()

		if i.Type == discordgo.InteractionApplicationCommand {
			if execution.IsDead(i.Member.User.ID) && !isHC(i.Member) {
				sender.RespondEphemeral(s, i, "https://tenor.com/view/yellow-emoji-no-no-emotiguy-no-no-no-gif-gif-9742000569423889376", nil)
			} else {
				log.Println("Received Command: " + i.ApplicationCommandData().Name)
				if h, ok := commands[i.ApplicationCommandData().Name]; ok {
					h.handler(s, i)
				}
			}
		}
	})
	s.ApplicationCommandBulkOverwrite(s.State.Application.ID, secrets.GuildID, commandSlice)
	fmt.Println("Done")

	fmt.Print("    |   Adding component commands... ")
	for commandName, commandFunction := range componentAndModalCommands {
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Type == discordgo.InteractionMessageComponent && strings.HasPrefix(strings.ToLower(i.MessageComponentData().CustomID), strings.ToLower(commandName)) {
				log.Println("Received Command: " + i.MessageComponentData().CustomID)
				if execution.IsDead(i.Member.User.ID) && !isHC(i.Member) {
					sender.RespondEphemeral(s, i, "https://tenor.com/view/yellow-emoji-no-no-emotiguy-no-no-no-gif-gif-9742000569423889376", nil)
				} else {
					commandFunction(s, i)
				}
			} else if i.Type == discordgo.InteractionModalSubmit && strings.HasPrefix(strings.ToLower(i.ModalSubmitData().CustomID), strings.ToLower(commandName)) {
				log.Println("Received Command: " + i.ModalSubmitData().CustomID)
				if execution.IsDead(i.Member.User.ID) && !isHC(i.Member) {
					sender.RespondEphemeral(s, i, "https://tenor.com/view/yellow-emoji-no-no-emotiguy-no-no-no-gif-gif-9742000569423889376", nil)
				} else {
					commandFunction(s, i)
				}
			}
		})
	}
	fmt.Println("Done")

	fmt.Print("    |   Adding autocompletions... ")
	for commandName, commandFunction := range autocompletes {
		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Type == discordgo.InteractionApplicationCommandAutocomplete && strings.HasPrefix(strings.ToLower(i.ApplicationCommandData().Name), strings.ToLower(commandName)) {
				log.Println("Received Autocomplete: " + i.ApplicationCommandData().Name)
				commandFunction(s, i)
			}
		})
	}
	fmt.Println("Done")
}

func isHC(m *discordgo.Member) bool {
	if IsAdminAbuser(m) {
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

func IsAdminAbuser(m *discordgo.Member) bool {
	return m.User.ID == "384422339393355786" || m.User.ID == "920342100468436993" || m.User.ID == "1079774043684745267" || m.User.ID == "952145898824138792"
}
