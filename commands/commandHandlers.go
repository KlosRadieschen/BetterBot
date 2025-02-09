package commands

import (
	"BetterScorch/execution"
	"BetterScorch/messages"
	"BetterScorch/sender"
	"fmt"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

func testHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Respond(s, i, "https://tenor.com/ss1MoenucUm.gif")
}

func executeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !execution.IsSacrificed(i.Member.User.ID) && !isHC(i.Member) {
		sender.RespondError(s, i, "The user is trying to execute a member, but they do not have the permissions to do that (they are a low ranking scum)")
	} else {
		sender.Respond(s, i, "Engaging target")
		execution.Execute(s, i.ApplicationCommandData().TargetID, i.ChannelID, false)
	}
}

func reviveHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	targetID := i.ApplicationCommandData().TargetID
	if !execution.IsDead(targetID) && !isHC(i.Member) {
		sender.RespondError(s, i, "The user is trying to revive an \"executed\" member, but the target is not even executed AND even if the target was downed, they do not even have the permissions to revive (they are a low ranking scum)")
	} else if !execution.IsDead(targetID) {
		sender.RespondError(s, i, "The user is trying to revive an \"executed\" member, but the target is not even executed")
	} else if !execution.IsSacrificed(targetID) && !isHC(i.Member) {
		sender.RespondError(s, i, "The user is trying to revive an executed member, but they do not have the permissions to do that (they are a low ranking scum)")
	} else {
		sender.Respond(s, i, "Commencing revive sequence")
		execution.Revive(s, i.ApplicationCommandData().TargetID, i.ChannelID)
	}
}

func toggleSleepHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messages.Sleeping = !messages.Sleeping
	if messages.Sleeping {
		sender.Respond(s, i, "https://tenor.com/view/ehouarn-sagot-dormir-a-mimir-mimir-gif-2358882822435654411")
	} else {
		sender.Respond(s, i, "https://tenor.com/view/wwe-coffin-world-wrestling-entertainment-gif-17903370")
	}
}

func statusHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if messages.Sleeping {
		sender.Respond(s, i, "https://tenor.com/view/dog-snoring-sleeping-meekotheiggy-knocked-out-gif-23834780")
	} else {
		sender.Respond(s, i, "https://tenor.com/view/funny-gif-15743464119256435424")
	}
}

func rollHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	max := 20
	reason := ""

	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "max":
			max = int(option.IntValue())
		case "reason":
			reason = fmt.Sprintf("Rolling for %v\n", option.StringValue())
		}
	}

	sender.Respond(s, i, fmt.Sprintf("%v%v/%v", reason, rand.Intn(max)+1, max))
}
