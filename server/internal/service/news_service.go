package service

import (
	"context"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/cache"
	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

type NewsService struct {
	newsRepo *repo.NewsRepo
	cache    *cache.Cache
}

func NewNewsService(newsRepo *repo.NewsRepo, cache *cache.Cache) *NewsService {
	return &NewsService{newsRepo: newsRepo, cache: cache}
}

func (s *NewsService) List(ctx context.Context, page, size int, category, section string) ([]*model.News, int, error) {
	return s.newsRepo.List(ctx, page, size, category, section)
}

func (s *NewsService) GetByID(ctx context.Context, id int64) (*model.News, error) {
	return s.newsRepo.FindByID(ctx, id)
}

func (s *NewsService) Search(ctx context.Context, q string, page, size int) ([]*model.News, int, error) {
	return s.newsRepo.Search(ctx, q, page, size)
}

func (s *NewsService) GetHeadlines(ctx context.Context) ([]*model.HeadlineWithNews, error) {
	cacheKey := "headlines:active"
	if val, ok := s.cache.Get(cacheKey); ok {
		return val.([]*model.HeadlineWithNews), nil
	}
	headlines, err := s.newsRepo.GetHeadlines(ctx)
	if err != nil {
		return nil, err
	}
	s.cache.SetWithTTL(cacheKey, headlines, 1, 2*time.Minute)
	return headlines, nil
}

func (s *NewsService) GetHeadlineHistory(ctx context.Context, date string) ([]*model.HeadlineWithNews, error) {
	return s.newsRepo.GetHeadlinesByDate(ctx, date)
}

func (s *NewsService) MarkRead(ctx context.Context, userID, newsID int64) error {
	return s.newsRepo.MarkRead(ctx, userID, newsID)
}

// Publish is called by the crawler after AI processing.
func (s *NewsService) Publish(ctx context.Context, n *model.News) (int64, error) {
	n.Status = 1
	now := time.Now()
	n.PublishedAt = &now
	return s.newsRepo.Create(ctx, n)
}

// SetHeadlines is used by admin to set headline news.
func (s *NewsService) SetHeadlines(ctx context.Context, headlines []*model.Headline) error {
	err := s.newsRepo.SetHeadlines(ctx, headlines)
	if err != nil {
		return err
	}
	s.cache.Del("headlines:active")
	return nil
}
