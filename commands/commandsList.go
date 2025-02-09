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
	{
		declaration: &discordgo.ApplicationCommand{
			Name:        "togglesleep",
			Description: "Toggle sleep on or off",
		},
		handler: toggleSleepHandler,
	},
	{
		declaration: &discordgo.ApplicationCommand{
			Name:        "sleepstatus",
			Description: "Find out if this fucker is awake",
		},
		handler: statusHandler,
	},
	{
		declaration: &discordgo.ApplicationCommand{
			Name:        "roll",
			Description: "Essentially gambling",
			Options: []*discordgo.ApplicationCommandOption{
				intOption("max", "The highest number the dice can get (default: 20)", false),
				stringOption("reason", "What you are rolling for", false),
			},
		},
		handler: rollHandler,
	},
}

func intOption(name string, desc string, required bool) *discordgo.ApplicationCommandOption {
	option := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        name,
		Description: desc,
		Required:    required,
	}
	return &option
}

func stringOption(name string, desc string, required bool) *discordgo.ApplicationCommandOption {
	option := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        name,
		Description: desc,
		Required:    required,
	}
	return &option
}

func boolOption(name string, desc string, required bool) *discordgo.ApplicationCommandOption {
	option := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        name,
		Description: desc,
		Required:    required,
	}
	return &option
}
