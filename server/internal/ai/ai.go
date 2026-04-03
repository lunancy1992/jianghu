package ai

import "context"

// AIClient is the interface for AI summarization and moderation.
type AIClient interface {
	Summarize(ctx context.Context, title, content string) (*SummaryResult, error)
	Moderate(ctx context.Context, content string) (*ModerationResult, error)
}
