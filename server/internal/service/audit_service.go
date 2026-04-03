package service

import (
	"context"

	"github.com/lunancy1992/jianghu-server/internal/ai"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

const (
	AuditPassed   = "passed"
	AuditRejected = "rejected"
	AuditPending  = "pending"
)

type AuditResult struct {
	Status string
	Reason string
}

type AuditService struct {
	aiClient    ai.AIClient
	commentRepo *repo.CommentRepo
}

func NewAuditService(aiClient ai.AIClient, commentRepo *repo.CommentRepo) *AuditService {
	return &AuditService{aiClient: aiClient, commentRepo: commentRepo}
}

// AuditContent calls AI moderator to check content.
func (s *AuditService) AuditContent(ctx context.Context, content string) (*AuditResult, error) {
	if s.aiClient == nil {
		// No AI client configured, default to pending
		return &AuditResult{Status: AuditPending}, nil
	}

	result, err := s.aiClient.Moderate(ctx, content)
	if err != nil {
		return &AuditResult{Status: AuditPending}, err
	}

	switch {
	case result.Passed:
		return &AuditResult{Status: AuditPassed}, nil
	case result.Rejected:
		return &AuditResult{Status: AuditRejected, Reason: result.Reason}, nil
	default:
		return &AuditResult{Status: AuditPending, Reason: result.Reason}, nil
	}
}

// GetPendingQueue returns pending comments for admin review.
func (s *AuditService) GetPendingQueue(ctx context.Context, page, size int) (interface{}, int, error) {
	return s.commentRepo.ListPending(ctx, page, size)
}

// Approve approves a pending comment.
func (s *AuditService) Approve(ctx context.Context, commentID int64) error {
	return s.commentRepo.UpdateStatus(ctx, commentID, 1, "admin approved")
}

// Reject rejects a pending comment.
func (s *AuditService) Reject(ctx context.Context, commentID int64, reason string) error {
	return s.commentRepo.UpdateStatus(ctx, commentID, 2, reason)
}
