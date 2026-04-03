package service

import (
	"context"

	"github.com/lunancy1992/jianghu-server/internal/model"
	"github.com/lunancy1992/jianghu-server/internal/repo"
)

type EventService struct {
	eventRepo *repo.EventRepo
}

func NewEventService(eventRepo *repo.EventRepo) *EventService {
	return &EventService{eventRepo: eventRepo}
}

func (s *EventService) List(ctx context.Context, page, size int) ([]*model.EventListItem, int, error) {
	return s.eventRepo.List(ctx, page, size)
}

type EventDetail struct {
	Event     *model.Event      `json:"event"`
	Nodes     []*model.EventNode `json:"nodes"`
	News      []*model.News      `json:"news"`
	Evidences []*model.Evidence  `json:"evidences"`
}

func (s *EventService) GetWithTimeline(ctx context.Context, id int64) (*EventDetail, error) {
	event, err := s.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, nil
	}

	nodes, err := s.eventRepo.ListNodes(ctx, id)
	if err != nil {
		return nil, err
	}

	news, err := s.eventRepo.GetLinkedNews(ctx, id)
	if err != nil {
		return nil, err
	}

	evidences, err := s.eventRepo.ListEvidences(ctx, id)
	if err != nil {
		return nil, err
	}

	return &EventDetail{
		Event:     event,
		Nodes:     nodes,
		News:      news,
		Evidences: evidences,
	}, nil
}

func (s *EventService) Create(ctx context.Context, e *model.Event) (int64, error) {
	if e.Status == "" {
		e.Status = "ongoing"
	}
	return s.eventRepo.Create(ctx, e)
}

func (s *EventService) AddNode(ctx context.Context, n *model.EventNode) (int64, error) {
	return s.eventRepo.CreateNode(ctx, n)
}

func (s *EventService) LinkNews(ctx context.Context, eventID, newsID int64) error {
	return s.eventRepo.LinkNews(ctx, eventID, newsID)
}

func (s *EventService) AddEvidence(ctx context.Context, e *model.Evidence) (int64, error) {
	exists, err := s.eventRepo.Exists(ctx, e.EventID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, nil
	}
	return s.eventRepo.CreateEvidence(ctx, e)
}
