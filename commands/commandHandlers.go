package commands

import (
	"BetterScorch/ai"
	"BetterScorch/execution"
	"BetterScorch/messages"
	"BetterScorch/polls"
	"BetterScorch/secrets"
	"BetterScorch/sender"
	"BetterScorch/webhooks"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func testHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Respond(s, i, "https://tenor.com/ss1MoenucUm.gif", nil)
}

func executeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !execution.IsSacrificed(i.Member.User.ID) && !isHC(i.Member) {
		sender.RespondError(s, i, "The user is trying to execute a member, but they do not have the permissions to do that (they are a low ranking scum)")
	} else {
		sender.Respond(s, i, "Engaging target", nil)

		var target string
		if i.ApplicationCommandData().TargetID == "" {
			target = i.ApplicationCommandData().Options[0].UserValue(nil).ID
		} else {
			target = i.ApplicationCommandData().TargetID
		}

		execution.Execute(s, target, i.ChannelID, false)
	}
}

func reviveHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var target string
	if i.ApplicationCommandData().TargetID == "" {
		target = i.ApplicationCommandData().Options[0].UserValue(nil).ID
	} else {
		target = i.ApplicationCommandData().TargetID
	}

	if !execution.IsDead(target) && !isHC(i.Member) {
		sender.RespondError(s, i, "The user is trying to revive an \"executed\" member, but the target is not even executed AND even if the target was downed, they do not even have the permissions to revive (they are a low ranking scum)")
	} else if !execution.IsDead(target) {
		sender.RespondError(s, i, "The user is trying to revive an \"executed\" member, but the target is not even executed")
	} else if !execution.IsSacrificed(target) && !isHC(i.Member) {
		sender.RespondError(s, i, "The user is trying to revive an executed member, but they do not have the permissions to do that (they are a low ranking scum)")
	} else {
		sender.Respond(s, i, "Commencing revive sequence", nil)
		execution.Revive(s, target, i.ChannelID)
	}
}

func pollHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Think(s, i)
	var pollOptions []string
	for _, option := range i.Interaction.ApplicationCommandData().Options[3:] {
		pollOptions = append(pollOptions, option.StringValue())
	}

	endTime := time.Now().Add(time.Duration(i.Interaction.ApplicationCommandData().Options[2].IntValue()) * time.Minute)

	pollID := polls.CreateOptionsPoll(s, i.Member.User.ID, i.Interaction.ApplicationCommandData().Options[1].BoolValue(), endTime, i.Interaction.ApplicationCommandData().Options[0].StringValue(), pollOptions...)
	sender.Followup(s, i, "Poll created")
	polls.WaitAndEvaluate(s, pollID, endTime)
}

func inputPollHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Think(s, i)

	endTime := time.Now().Add(time.Duration(i.Interaction.ApplicationCommandData().Options[2].IntValue()) * time.Minute)

	pollID := polls.CreateInputPoll(s, i.Member.User.ID, i.Interaction.ApplicationCommandData().Options[1].BoolValue(), endTime, i.Interaction.ApplicationCommandData().Options[0].StringValue())
	sender.Followup(s, i, "Poll created")
	polls.WaitAndEvaluateInput(s, pollID, endTime)
}

func toggleSleepHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	messages.Sleeping = !messages.Sleeping
	if messages.Sleeping {
		sender.Respond(s, i, "https://tenor.com/view/ehouarn-sagot-dormir-a-mimir-mimir-gif-2358882822435654411", nil)
	} else {
		sender.Respond(s, i, "https://tenor.com/view/wwe-coffin-world-wrestling-entertainment-gif-17903370", nil)
	}
}

func statusHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if messages.Sleeping {
		sender.Respond(s, i, "https://tenor.com/view/dog-snoring-sleeping-meekotheiggy-knocked-out-gif-23834780", nil)
	} else {
		sender.Respond(s, i, "https://tenor.com/view/funny-gif-15743464119256435424", nil)
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

	sender.Respond(s, i, fmt.Sprintf("%v%v/%v", reason, rand.Intn(max)+1, max), nil)
}

func addPersonalityHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Think(s, i)
	nickname := i.ApplicationCommandData().Options[0].StringValue()
	pfpLink := ""
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case "nickname":
			nickname = option.StringValue()
		case "pfplink":
			pfpLink = option.StringValue()
		}
	}

	name := i.ApplicationCommandData().Options[0].StringValue()
	appropriate, err := webhooks.IsAppropriate(name)
	if sender.HandleErrInteractionFollowup(s, i, err) {
		return
	} else if appropriate && len(name) <= 80 && !strings.Contains(strings.ToLower(name), "discord") && !strings.Contains(strings.ToLower(name), "clyde") {
		webhooks.AddPersonality(s, i, name, nickname, pfpLink)
		sender.Followup(s, i, fmt.Sprintf("%v joined the chat", nickname))
	} else {
		sender.RespondError(s, i, "They tried to add an AI personality but the name \""+name+"\" was deemed inappropriate")
	}
}

func killPersonalityHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	name := i.ApplicationCommandData().Options[0].StringValue()
	if webhooks.PersonalityExists(name) {
		sender.Respond(s, i, "I'm shooting "+name, nil)
		webhooks.RemovePersonality(s, i, name)
	} else {
		sender.RespondError(s, i, "The user is trying to remove the AI personality \""+name+"\" but it does not exist")
	}
}

func purgeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if webhooks.PersonalitiesEmpty() {
		sender.RespondError(s, i, "The user is trying to purge AI personalities but there are currently none")
	} else {
		sender.Respond(s, i, "https://tenor.com/view/langley-thanos-gif-20432464", nil)
		webhooks.Purge(s, i)
	}
}

func exposeHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if messages.Msgs[i.ApplicationCommandData().Options[0].UserValue(nil).ID] == nil {
		sender.Respond(s, i, "User doesn't have any messages", nil)
		return
	}

	member, _ := s.GuildMember(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID)
	embed := discordgo.MessageEmbed{
		Title: "Exposing " + member.Nick,
		Color: 0xFF69B4,
	}
	for msg := messages.Msgs[i.ApplicationCommandData().Options[0].UserValue(nil).ID].Front(); msg != nil; msg = msg.Next() {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("<t:%v:R>", msg.Value.(*discordgo.Message).Timestamp.Unix()),
			Value:  msg.Value.(*discordgo.Message).Content,
			Inline: false,
		})
	}
	sender.Respond(s, i, "", []*discordgo.MessageEmbed{&embed})
}

func linkHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	sender.Think(s, i)
	resp, err := ai.GenerateSingleResponse("Write an extremely annoying and obnoxious short (only one paragraph) ad for the AHA website which lets you read reports from the AHA and modify your character. Put in the link like this: [AHA website](https://aha-rp.org)")
	if !sender.HandleErrInteractionFollowup(s, i, err) {
		sender.Followup(s, i, resp)
	}
}

func promoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !isHC(i.Member) {
		sender.RespondError(s, i, "User tried to promote someone but doesn't have the proper perms")
		return
	}

	var roles = []string{"1195135956471255140", "1195858311627669524", "1195858271349784639", "1195136106811887718", "1195858179590987866", "1195137362259349504", "1195136284478410926", "1195137253408768040", "1250582641921757335", "1195758308519325716", "1195758241221722232", "1195758137563689070", "1195757362439528549", "1195136491148550246", "1195708423229165578", "1195137477497868458", "1195136604373782658", "1248776818664935525", "1277091839647678575"}

	count := 1
	if len(i.ApplicationCommandData().Options) == 3 {
		count = int(i.ApplicationCommandData().Options[2].IntValue())
	}

	member, _ := s.GuildMember(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID)
	for _, role := range member.Roles {
		if slices.Contains(roles, role) {
			s.GuildMemberRoleRemove(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID, role)
			s.GuildMemberRoleAdd(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID, roles[slices.Index(roles, role)-count])
		}
	}

	sender.Respond(s, i, member.Mention()+", you have been promoted:\n"+i.ApplicationCommandData().Options[1].StringValue(), nil)
}

func demoteHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if !isHC(i.Member) {
		sender.RespondError(s, i, "User tried to promote someone but doesn't have the proper perms")
		return
	}

	var roles = []string{"1195135956471255140", "1195858311627669524", "1195858271349784639", "1195136106811887718", "1195858179590987866", "1195137362259349504", "1195136284478410926", "1195137253408768040", "1250582641921757335", "1195758308519325716", "1195758241221722232", "1195758137563689070", "1195757362439528549", "1195136491148550246", "1195708423229165578", "1195137477497868458", "1195136604373782658", "1248776818664935525", "1277091839647678575"}

	count := 1
	if len(i.ApplicationCommandData().Options) == 3 {
		count = int(i.ApplicationCommandData().Options[2].IntValue())
	}

	member, _ := s.GuildMember(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID)
	for _, role := range member.Roles {
		if slices.Contains(roles, role) {
			s.GuildMemberRoleRemove(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID, role)
			s.GuildMemberRoleAdd(secrets.GuildID, i.ApplicationCommandData().Options[0].UserValue(nil).ID, roles[slices.Index(roles, role)+count])
		}
	}

	sender.Respond(s, i, member.Mention()+", you have been demoted:\n"+i.ApplicationCommandData().Options[1].StringValue(), nil)
}

func messageHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	charName := ""
	if len(i.ApplicationCommandData().Options) > 1 {
		exists := false
		for ID, characters := range webhooks.CharacterBuffer {
			if ID == i.Member.User.ID {
				for _, character := range characters {
					if character.Name == i.ApplicationCommandData().Options[1].StringValue() {
						exists = true
					}
				}
			}
		}

		if exists {
			charName = "-" + i.ApplicationCommandData().Options[1].StringValue()
		} else {
			sender.RespondError(s, i, "User is trying to send a message using a character that doesn't exist")
			return
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Send message",
			CustomID: "messagemodalsubmit-" + i.ApplicationCommandData().Options[0].UserValue(nil).ID + charName,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID: "message",
							Label:    "Message",
							Style:    discordgo.TextInputParagraph,
							Required: true,
						},
					},
				},
			},
		},
	})
	sender.HandleErrInteraction(s, i, err)
}
