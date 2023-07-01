package utilities

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RegisterCommand(s *discordgo.Session, event *discordgo.GuildCreate) {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "create-channel",
			Description: "Create a channel to post youtube videos of the specified search words.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channel_name",
					Description: "Write the name of the channel",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "search_words",
					Description: "Write the search words",
					Required:    true,
				},
			},
		},
		{
			Name:        "register-apikey",
			Description: "Register the apikey to post youtube videos",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "apikey",
					Description: "Write the apikey of Youtube",
					Required:    true,
				},
			},
		},
		{
			Name:        "add-blacklist",
			Description: "Add the channel to the blacklist[format:youtube.com/@~~~]",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channel-url",
					Description: "Please post the Youtube channel link you wish to add to the blacklist",
					Required:    true,
				},
			},
		},
		{
			Name:        "remove-blacklist",
			Description: "Remove the channel from the blacklist[format:youtube.com/@~~~]",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "channel-url",
					Description: "Please post the Youtube channel link you wish to remove from the blacklist",
					Required:    true,
				},
			},
		},
		{
			Name:        "get-blacklist",
			Description: "Get the blacklist",
		},
		{
			Name:        "start-notification",
			Description: "Start notifications of youtube search results at specified time intervals (unit: hours)",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "time_interval",
					Description: "Integer from 1 to 23.Default 3",
					Required:    false,
				},
			},
		},
		{
			Name:        "stop-notification",
			Description: "Stop notifications of youtube search results",
		},
		{
			Name:        "help",
			Description: "Get the command list",
		},
	}

	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, event.Guild.ID, command)
		if err != nil {
			log.Println("Failed to create command", command.Name, err)
		}
	}
}
