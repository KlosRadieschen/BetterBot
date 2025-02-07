package commands

import (
	"github.com/bwmarrin/discordgo"
)

var commands = []command{
	{
		declaration: &discordgo.ApplicationCommand{
			Name:        "test",
			Description: "Test if this fucker is online",
		},
		handler: testHandler,
	},
	{
		declaration: &discordgo.ApplicationCommand{
			Type: discordgo.UserApplicationCommand,
			Name: "FUCKING KILL THEM",
		},
		handler: executeHandler,
	},
	{
		declaration: &discordgo.ApplicationCommand{
			Type: discordgo.UserApplicationCommand,
			Name: "Yeah I guess they can live",
		},
		handler: reviveHandler,
	},
}
