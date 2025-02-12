package sender

import (
	"BetterScorch/ai"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func SetResponseTimeout(s *discordgo.Session, i *discordgo.InteractionCreate, duration time.Duration) {
	time.Sleep(duration)
	s.InteractionResponseDelete(i.Interaction)
}

func Respond(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: resp,
		},
	})
}

func RespondEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{

		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: resp,
		},
	})
}

func Followup(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: resp,
	})
}

func FollowupEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate, resp string) *discordgo.InteractionCreate {
	s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Flags:   discordgo.MessageFlagsEphemeral,
		Content: resp,
	})
	return i
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

func ThinkEphemeral(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}

func HandleErrInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, err error) bool {
	if err != nil {
		log.Printf("Received error: %s", err.Error())
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error:\n```" + err.Error() + "```",
			},
		})
		return true
	}
	return false
}

func HandleErrInteractionFollowup(s *discordgo.Session, i *discordgo.InteractionCreate, err error) bool {
	if err != nil {
		log.Printf("Received error: %s", err.Error())
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Error:\n```" + err.Error() + "```",
		})
		return true
	}
	return false
}
