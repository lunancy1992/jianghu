package service

import (
	"context"
	"fmt"
	"log"

	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

type CommentService struct {
	commentRepo  *repo.CommentRepo
	coinRepo     *repo.CoinRepo
	newsRepo     *repo.NewsRepo
	auditService *AuditService
}

func NewCommentService(commentRepo *repo.CommentRepo, coinRepo *repo.CoinRepo, newsRepo *repo.NewsRepo, auditService *AuditService) *CommentService {
	return &CommentService{
		commentRepo:  commentRepo,
		coinRepo:     coinRepo,
		newsRepo:     newsRepo,
		auditService: auditService,
	}
}

// Create handles comment creation: AI audit + deduct coin + insert.
func (s *CommentService) Create(ctx context.Context, newsID, userID int64, content, stance string) (*model.Comment, error) {
	news, err := s.newsRepo.FindByID(ctx, newsID)
	if err != nil {
		return nil, err
	}
	if news == nil {
		return nil, fmt.Errorf("news not found")
	}

	// Step 1: AI audit
	auditResult, err := s.auditService.AuditContent(ctx, content)
	if err != nil {
		log.Printf("AI audit error, defaulting to pending: %v", err)
		auditResult = &AuditResult{Status: AuditPending}
	}

	if auditResult.Status == AuditRejected {
		return nil, fmt.Errorf("comment rejected: %s", auditResult.Reason)
	}

	// Step 2: Deduct coin (1 coin per comment)
	tx, err := s.coinRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = s.coinRepo.Deduct(ctx, tx, userID, 1, "comment_cost", "", "发表评论消耗")
	if err != nil {
		return nil, fmt.Errorf("insufficient coins to comment")
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Step 3: Insert comment
	status := 1 // approved
	if auditResult.Status == AuditPending {
		status = 0 // pending review
	}

	comment := &model.Comment{
		NewsID:  newsID,
		UserID:  userID,
		Content: content,
		Stance:  stance,
		Status:  status,
	}

	id, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}

	// Re-fetch from DB to get complete fields (created_at, updated_at)
	created, err := s.commentRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (s *CommentService) List(ctx context.Context, newsID int64, page, size int) ([]*model.Comment, int, error) {
	return s.commentRepo.ListByNews(ctx, newsID, page, size)
}

func (s *CommentService) Like(ctx context.Context, userID, commentID int64) error {
	return s.commentRepo.CreateLike(ctx, userID, commentID)
}

func (s *CommentService) Unlike(ctx context.Context, userID, commentID int64) error {
	return s.commentRepo.DeleteLike(ctx, userID, commentID)
}
