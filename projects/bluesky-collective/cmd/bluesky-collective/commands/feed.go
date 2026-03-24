package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/consensuscode/bluesky-collective/pkg/atproto"
)

// NewFeedCmd creates the feed command for reading Bluesky feeds.
func NewFeedCmd(logger *zap.Logger) *cobra.Command {
	var (
		actor string
		limit int
	)

	cmd := &cobra.Command{
		Use:   "feed",
		Short: "Read a Bluesky user's feed",
		Long: `Feed retrieves recent posts from a Bluesky user's feed.

This is a read-only operation that does not require consensus.

Example:
  bluesky-collective feed --actor collectiveflow.bsky.social
  bluesky-collective feed --actor collectiveflow.bsky.social --limit 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if actor == "" {
				actor = viper.GetString("bluesky.handle")
			}
			if actor == "" {
				return fmt.Errorf("actor is required: use --actor or configure bluesky.handle")
			}

			logger.Info("Fetching feed",
				zap.String("actor", actor),
				zap.Int("limit", limit),
			)

			serviceURL := viper.GetString("bluesky.service")
			if serviceURL == "" {
				serviceURL = "https://bsky.social"
			}

			client := atproto.NewClient(serviceURL)

			// Authenticate if credentials are available (some feeds are public)
			identifier := viper.GetString("bluesky.identifier")
			password := viper.GetString("bluesky.password")
			if identifier != "" && password != "" {
				if err := client.Authenticate(cmd.Context(), identifier, password); err != nil {
					logger.Warn("Authentication failed, trying unauthenticated",
						zap.Error(err),
					)
				}
			}

			feed, err := client.GetAuthorFeed(cmd.Context(), actor, limit)
			if err != nil {
				return fmt.Errorf("get feed for %s: %w", actor, err)
			}

			if len(feed.Feed) == 0 {
				fmt.Printf("No posts found for @%s\n", actor)
				return nil
			}

			fmt.Printf("Feed for @%s (%d posts):\n\n", actor, len(feed.Feed))
			for i, item := range feed.Feed {
				// Parse the record to extract text
				var record struct {
					Text      string `json:"text"`
					CreatedAt string `json:"createdAt"`
				}
				if err := json.Unmarshal(item.Post.Record, &record); err != nil {
					record.Text = "(unable to parse)"
				}

				fmt.Printf("%d. @%s", i+1, item.Post.Author.Handle)
				if item.Post.Author.DisplayName != "" {
					fmt.Printf(" (%s)", item.Post.Author.DisplayName)
				}
				fmt.Println()
				fmt.Printf("   %s\n", record.Text)
				fmt.Printf("   %s\n", record.CreatedAt)
				fmt.Printf("   URI: %s\n\n", item.Post.URI)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&actor, "actor", "a", "", "Bluesky handle to fetch feed for (defaults to configured handle)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of posts to fetch (1-100)")

	return cmd
}
