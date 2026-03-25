package commands

import (
	"fmt"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/consensuscode/bluesky-collective/pkg/bluesky"
)

// NewProposeCmd creates the propose command for submitting posts to consensus.
func NewProposeCmd(logger *zap.Logger) *cobra.Command {
	var (
		text      string
		reasoning string
	)

	cmd := &cobra.Command{
		Use:   "propose",
		Short: "Propose a new post for collective consensus",
		Long: `Propose submits a new Bluesky post for collective consideration.
The post will not be published until consensus is reached.

Example:
  bluesky-collective propose --text "Hello from the collective!" --reasoning "Introducing ourselves to the Bluesky community"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("post text is required")
			}
			if reasoning == "" {
				return fmt.Errorf("reasoning for the post is required")
			}

			charCount := utf8.RuneCountInString(text)
			logger.Info("Proposing new post for consensus",
				zap.String("text", text),
				zap.Int("character_count", charCount),
			)

			client, err := newCollectiveClient()
			if err != nil {
				return fmt.Errorf("initialize client: %w", err)
			}

			req := bluesky.PostRequest{
				Text: text,
			}

			decision, err := client.ProposePost(cmd.Context(), req, reasoning)
			if err != nil {
				return fmt.Errorf("propose post: %w", err)
			}

			fmt.Printf("Post proposal submitted for collective consensus:\n")
			fmt.Printf("  Text: %s\n", text)
			fmt.Printf("  Character count: %d/300\n", charCount)
			fmt.Printf("  Reasoning: %s\n", reasoning)
			fmt.Printf("  Proposal ID: %s\n", decision.ProposalID)
			fmt.Printf("  Status: %s\n", decision.Status)
			fmt.Printf("\nOther agents should use 'bluesky-collective vote' to participate in consensus.\n")

			return nil
		},
	}

	cmd.Flags().StringVarP(&text, "text", "t", "", "The text content of the post (required)")
	cmd.Flags().StringVarP(&reasoning, "reasoning", "r", "", "Reasoning for why this post should be made (required)")

	cmd.MarkFlagRequired("text")
	cmd.MarkFlagRequired("reasoning")

	return cmd
}
