package commands

import (
	"github.com/bwmarrin/discordgo"
)

var commands = map[string]command{
	/*

		Execute

	*/

	"FUCKING KILL THEM": {
		declaration: &discordgo.ApplicationCommand{
			Type: discordgo.UserApplicationCommand,
			Name: "FUCKING KILL THEM",
		},
		handler: executeHandler,
	},
	"Yeah I guess they can live": {
		declaration: &discordgo.ApplicationCommand{
			Type: discordgo.UserApplicationCommand,
			Name: "Yeah I guess they can live",
		},
		handler: reviveHandler,
	},
	"execute": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "execute",
			Description: "Admin abuse my beloved",
			Options: []*discordgo.ApplicationCommandOption{
				userOption("target", "User you want to kill", true),
			},
		},
		handler: executeHandler,
	},
	"revive": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "revive",
			Description: "Admin abuse my beloved",
			Options: []*discordgo.ApplicationCommandOption{
				userOption("target", "USer you want to revive", true),
			},
		},
		handler: reviveHandler,
	},

	/*

		Polls

	*/

	"poll": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "poll",
			Description: "Ask the people of the AHA",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "question",
					Description: "The question you want to ask",
					Required:    true,
					MaxLength:   120,
				},
				boolOption("multi-option", "Whether or not users can vote for multiple options", true),
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Duration of the poll",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "1 minute",
							Value: 1,
						},
						{
							Name:  "15 minutes",
							Value: 15,
						},
						{
							Name:  "1 hour",
							Value: 60,
						},
						{
							Name:  "6 hours",
							Value: 6 * 60,
						},
						{
							Name:  "12 hours",
							Value: 12 * 60,
						},
						{
							Name:  "1 day",
							Value: 24 * 60,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option1",
					Description: "First option that people can choose",
					Required:    true,
					MaxLength:   25,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option2",
					Description: "Second option that people can choose",
					Required:    true,
					MaxLength:   25,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option3",
					Description: "Third option that people can choose",
					Required:    false,
					MaxLength:   25,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option4",
					Description: "Fourth option that people can choose",
					Required:    false,
					MaxLength:   25,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "option5",
					Description: "Fifth option that people can choose",
					Required:    false,
					MaxLength:   25,
				},
			},
		},
		handler: pollHandler,
	},
	"inputpoll": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "inputpoll",
			Description: "Ask the people of the AHA",
			Options: []*discordgo.ApplicationCommandOption{
				stringOption("question", "The question you want to ask", true),
				boolOption("multioption", "Whether or not users can submit multiple messages", true),
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "Duration of the poll",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "1 minute",
							Value: 1,
						},
						{
							Name:  "15 minutes",
							Value: 15,
						},
						{
							Name:  "1 hour",
							Value: 60,
						},
						{
							Name:  "6 hours",
							Value: 6 * 60,
						},
						{
							Name:  "12 hours",
							Value: 12 * 60,
						},
						{
							Name:  "1 day",
							Value: 24 * 60,
						},
					},
				},
			},
		},
		handler: inputPollHandler,
	},

	/*

		Sleep

	*/

	"togglesleep": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "togglesleep",
			Description: "Toggle sleep on or off",
		},
		handler: toggleSleepHandler,
	},
	"sleepstatus": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "sleepstatus",
			Description: "Find out if this fucker is awake",
		},
		handler: statusHandler,
	},

	/*

		Personality

	*/

	"addpersonality": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "addpersonality",
			Description: "Add an AI personality to the server",
			Options: []*discordgo.ApplicationCommandOption{
				stringOption("name", "The FULL name of the personality", true),
				stringOption("nickname", "The name that will trigger the personality", false),
				stringOption("pfplink", "A link to an image which will be used as PFP", false),
			},
		},
		handler: addPersonalityHandler,
	},
	"killpersonality": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "killpersonality",
			Description: "Remove an AI personality from the server",
			Options: []*discordgo.ApplicationCommandOption{
				stringOption("name", "The FULL name of the personality", true),
			},
		},
		handler: killPersonalityHandler,
	},
	"purgepersonalities": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "purgepersonalities",
			Description: "Remove ALL AI personality from the server",
		},
		handler: purgeHandler,
	},

	/*

		Characters

	*/

	"addcharacter": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "addcharacter",
			Description: "Add a new tupper-like character",
			Options: []*discordgo.ApplicationCommandOption{
				stringOption("name", "Name of the character", true),
				stringOption("brackets", "Define the trigger for the character", true),
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "pfp",
					Description: "Profile picture of the character",
					Required:    true,
				},
			},
		},
		handler: addCharacterHandler,
	},
	"removecharacter": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "removecharacter",
			Description: "Remove one of yourcharacter",
			Options: []*discordgo.ApplicationCommandOption{
				stringOption("name", "Name of the character", true),
			},
		},
		handler: removeCharacterHandler,
	},
	"listcharacters": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "listcharacters",
			Description: "List all your characters",
		},
		handler: listCharactersHandler,
	},

	/*

		Register

	*/

	"register": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "register",
			Description: "Register yourself into the database",
			Options: []*discordgo.ApplicationCommandOption{
				stringOption("full-name", " Full name of your character", true),
				stringOption("callsign", "Callsign of your character", true),
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "picture",
					Description: "Picture of the character",
					Required:    true,
				},
			},
		},
		handler: registerHandler,
	},
	"unregister": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "unregister",
			Description: "Unregister yourself from the database",
		},
		handler: unregisterHandler,
	},

	/*

		Reports

	*/

	"listreports": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "listreports",
			Description: "List all reports",
		},
		handler: listReportsHandler,
	},
	"getreport": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "getreport",
			Description: "Get a report",
			Options: []*discordgo.ApplicationCommandOption{
				intOption("index", "Index of the report", true),
			},
		},
		handler: getReportHandler,
	},

	/*

		Miscellaneous

	*/

	"promote": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "promote",
			Description: "Please stop bitching Kerm",
			Options: []*discordgo.ApplicationCommandOption{
				userOption("user", "Person receiving the promotion", true),
				stringOption("reason", "Reason for the promotion", true),
				intOption("amount", "Amount of promotions", false),
			},
		},
		handler: promoteHandler,
	},
	"demote": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "demote",
			Description: "Please stop bitching Kerm",
			Options: []*discordgo.ApplicationCommandOption{
				userOption("user", "Person receiving the demotion", true),
				stringOption("reason", "Reason for the demotion", true),
				intOption("amount", "Amount of demotions", false),
			},
		},
		handler: demoteHandler,
	},

	/*

		Miscellaneous

	*/

	"test": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "test",
			Description: "Test if this fucker is online",
		},
		handler: testHandler,
	},
	"roll": {
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
	"expose": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "expose",
			Description: "Show the target's recent messages",
			Options: []*discordgo.ApplicationCommandOption{
				userOption("target", "The target that you want to expose", true),
			},
		},
		handler: exposeHandler,
	},
	"link": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "link",
			Description: "Get the link to the website",
		},
		handler: linkHandler,
	},
	"message": {
		declaration: &discordgo.ApplicationCommand{
			Name:        "message",
			Description: "Leave a message for someone",
			Options: []*discordgo.ApplicationCommandOption{
				userOption("recipient", "Person receiving the message", true),
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "character",
					Description:  "Sends the message with one of your characters",
					Required:     false,
					Autocomplete: true,
				},
			},
		},
		handler: messageHandler,
	},
}

var componentAndModalCommands = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"pollvote":             pollVoteHandler,
	"pollshow":             pollShowHandler,
	"inputpollvote":        inputPollVoteHandler,
	"inputpollmodalcreate": inputPollModalCreateHandler,
	"inputpollmodalsubmit": inputPollModalSubmit,
	"inputpollshow":        inputPollShowHandler,
	"messagemodalsubmit":   messageModalSubmit,
}

var autocompletes = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"message": messageAutocompleteHandler,
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

func userOption(name string, desc string, required bool) *discordgo.ApplicationCommandOption {
	option := discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        name,
		Description: desc,
		Required:    required,
	}
	return &option
}
