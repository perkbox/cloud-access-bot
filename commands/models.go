package commands

import (
	"github.com/perkbox/cloud-access-bot/internal"
	"github.com/perkbox/cloud-access-bot/internal/settings"
	"github.com/slack-go/slack/socketmode"
)

// SlashCommandController We create a structure to let us use dependency injection
type SlashCommandController struct {
	EventHandler *socketmode.SocketmodeHandler
	Service      internal.Service
	Settings     settings.Settings
}
