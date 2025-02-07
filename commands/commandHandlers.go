package commands

import (
	"BetterScorch/ai"
	"BetterScorch/execution"

	"github.com/bwmarrin/discordgo"
)

func testHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	respond(s, i, "https://tenor.com/ss1MoenucUm.gif")
}

func executeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !execution.IsSacrificed(i.Member.User.ID) && !isHC(i.Member) {
		respondError(s, i, "The user is trying to execute a member, but they do not have the permissions to do that (they are a low ranking scum)")
	} else {
		respond(s, i, "Engaging target")
		execution.Execute(s, i.ApplicationCommandData().TargetID, i.ChannelID, false)
	}
}

func reviveHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	targetID := i.ApplicationCommandData().TargetID
	if !execution.IsDead(targetID) && !isHC(i.Member) {
		respondError(s, i, "The user is trying to revive an \"executed\" member, but the target is not even executed AND even if the target was downed, they do not even have the permissions to revive (they are a low ranking scum)")
	} else if !execution.IsDead(targetID) {
		respondError(s, i, "The user is trying to revive an \"executed\" member, but the target is not even executed")
	} else if !execution.IsSacrificed(targetID) && !isHC(i.Member) {
		respondError(s, i, "The user is trying to revive an executed member, but they do not have the permissions to do that (they are a low ranking scum)")
	} else {
		respond(s, i, "Commencing revive sequence")
		execution.Revive(s, i.ApplicationCommandData().TargetID, i.ChannelID)
	}
}

func respond(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: resp,
		},
	})
}

func followup(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: resp,
	})
}

func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, errorDescription string) {
	think(s, i)
	errorResponse, err := ai.GenerateErrorResponse(errorDescription)
	handleErrFollowup(s, i, err)
	followup(s, i, errorResponse)
}

func think(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func handleErr(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error:\n```" + err.Error() + "```",
			},
		})
	}
}

func handleErrFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Error:\n```" + err.Error() + "```",
		})
	}
}
