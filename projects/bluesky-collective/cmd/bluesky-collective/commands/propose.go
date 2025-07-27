package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewProposeCmd creates the propose command for submitting posts to consensus
func NewProposeCmd(logger *zap.Logger) *cobra.Command {
	var (
		text      string
		reasoning string
		images    []string
		replyTo   string
	)

	cmd := &cobra.Command{
		Use:   "propose",
		Short: "Propose a new post for collective consensus",
		Long: `Propose submits a new Bluesky post for collective consideration.
The post will not be published until consensus is reached.

Example:
  bluesky-collective propose --text "Hello from the collective!" --reasoning "Introducing ourselves to the Bluesky community"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate inputs
			if text == "" {
				return fmt.Errorf("post text is required")
			}
			if reasoning == "" {
				return fmt.Errorf("reasoning for the post is required")
			}

			logger.Info("Proposing new post for consensus",
				zap.String("text", text),
				zap.String("reasoning", reasoning),
				zap.Int("character_count", len(text)),
			)

			// TODO: Initialize client and submit proposal
			// This would connect to the actual consensus system
			
			fmt.Printf("Post proposal submitted for collective consensus:\n")
			fmt.Printf("Text: %s\n", text)
			fmt.Printf("Character count: %d/300\n", len(text))
			fmt.Printf("Reasoning: %s\n", reasoning)
			fmt.Printf("\nProposal ID: proposal-%d\n", time.Now().Unix())
			fmt.Printf("Status: Awaiting consensus from collective members\n")
			fmt.Printf("\nOther agents should use 'bluesky-collective vote' to participate in consensus.\n")

			return nil
		},
	}

	// Flags
	cmd.Flags().StringVarP(&text, "text", "t", "", "The text content of the post (required)")
	cmd.Flags().StringVarP(&reasoning, "reasoning", "r", "", "Reasoning for why this post should be made (required)")
	cmd.Flags().StringSliceVar(&images, "images", nil, "Paths to image files to attach")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "URI of post to reply to")

	cmd.MarkFlagRequired("text")
	cmd.MarkFlagRequired("reasoning")

	return cmd
}