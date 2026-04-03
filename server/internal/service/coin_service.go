package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

type CoinService struct {
	coinRepo    *repo.CoinRepo
	commentRepo *repo.CommentRepo
}

func NewCoinService(coinRepo *repo.CoinRepo, commentRepo *repo.CommentRepo) *CoinService {
	return &CoinService{coinRepo: coinRepo, commentRepo: commentRepo}
}

func (s *CoinService) GetBalance(ctx context.Context, userID int64) (int64, error) {
	acc, err := s.coinRepo.GetAccount(ctx, userID)
	if err != nil {
		return 0, err
	}
	if acc == nil {
		return 0, nil
	}
	return acc.Balance, nil
}

func (s *CoinService) GetTransactions(ctx context.Context, userID int64, page, size int) ([]*model.CoinTransaction, int, error) {
	return s.coinRepo.ListTransactions(ctx, userID, page, size)
}

// WeeklyDistribute grants coins to all active users (called by cron Monday 00:00).
func (s *CoinService) WeeklyDistribute(ctx context.Context) error {
	since := time.Now().AddDate(0, 0, -7)
	activeUsers, err := s.coinRepo.GetActiveUsers(ctx, since)
	if err != nil {
		return fmt.Errorf("get active users: %w", err)
	}

	log.Printf("[CoinService] Weekly distribution to %d active users", len(activeUsers))

	for _, uid := range activeUsers {
		tx, err := s.coinRepo.BeginTx(ctx)
		if err != nil {
			log.Printf("[CoinService] begin tx for user %d: %v", uid, err)
			continue
		}

		err = s.coinRepo.Credit(ctx, tx, uid, 3, "weekly_grant", "", "每周活跃奖励")
		if err != nil {
			tx.Rollback()
			log.Printf("[CoinService] credit user %d: %v", uid, err)
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("[CoinService] commit user %d: %v", uid, err)
		}
	}

	return nil
}

// ProcessCommentRewards rewards top-liked comments (called hourly by cron).
func (s *CoinService) ProcessCommentRewards(ctx context.Context) error {
	since := time.Now().Add(-2 * time.Hour) // look at last 2 hours
	topComments, err := s.commentRepo.FindTopLiked(ctx, since, 10)
	if err != nil {
		return err
	}

	for _, c := range topComments {
		tx, err := s.coinRepo.BeginTx(ctx)
		if err != nil {
			continue
		}

		refID := fmt.Sprintf("comment:%d", c.ID)
		err = s.coinRepo.Credit(ctx, tx, c.UserID, 2, "comment_reward", refID, "热评奖励")
		if err != nil {
			tx.Rollback()
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("[CoinService] reward commit for comment %d: %v", c.ID, err)
		}
	}

	return nil
}
