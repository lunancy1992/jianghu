package cron

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	c *cron.Cron
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		c: cron.New(cron.WithSeconds()),
	}
}

func (s *Scheduler) AddFunc(spec string, cmd func()) error {
	_, err := s.c.AddFunc(spec, cmd)
	if err != nil {
		log.Printf("[Cron] Failed to add job %q: %v", spec, err)
		return err
	}
	log.Printf("[Cron] Registered job: %s", spec)
	return nil
}

func (s *Scheduler) Start() {
	log.Println("[Cron] Starting scheduler")
	s.c.Start()
}

func (s *Scheduler) Stop() {
	log.Println("[Cron] Stopping scheduler")
	ctx := s.c.Stop()
	<-ctx.Done()
}
