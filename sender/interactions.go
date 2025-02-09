package sender

import (
	"BetterScorch/ai"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Respond(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: resp,
		},
	})
}

func Followup(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: resp,
	})
}

func RespondError(s *discordgo.Session, i *discordgo.InteractionCreate, errorDescription string) {
	Think(s, i)
	errorResponse, err := ai.GenerateErrorResponse(errorDescription)
	HandleErrInteractionFollowup(s, i, err)
	Followup(s, i, errorResponse)
}

func Think(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func HandleErrInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	log.Printf("Received error: %s", err.Error())
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error:\n```" + err.Error() + "```",
			},
		})
	}
}

func HandleErrInteractionFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Error:\n```" + err.Error() + "```",
		})
	}
}
