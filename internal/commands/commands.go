package commands

import (
	"github.com/musicmash/musicmash/internal/commands/artists"
	"github.com/musicmash/musicmash/internal/commands/releases"
	"github.com/musicmash/musicmash/internal/commands/stores"
	"github.com/musicmash/musicmash/internal/commands/subscriptions"
	"github.com/spf13/cobra"
)

func AddCommands(cmd *cobra.Command) {
	cmd.AddCommand(
		artists.NewArtistCommand(),
		stores.NewStoreCommand(),
		subscriptions.NewSubscriptionCommand(),
		releases.NewReleaseCommand(),
	)
}
