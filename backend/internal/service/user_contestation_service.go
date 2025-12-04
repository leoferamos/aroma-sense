package service

import (
	"time"

	"github.com/leoferamos/aroma-sense/internal/model"
	"github.com/leoferamos/aroma-sense/internal/repository"
)

type UserContestationService interface {
	Create(userID uint, reason string) error
	ListPending(limit, offset int) ([]model.UserContestation, int64, error)
	Approve(id uint, adminID uint, notes *string) error
	Reject(id uint, adminID uint, notes *string) error
}

type userContestationService struct {
	repo repository.UserContestationRepository
}

func NewUserContestationService(repo repository.UserContestationRepository) UserContestationService {
	return &userContestationService{repo: repo}
}

func (s *userContestationService) Create(userID uint, reason string) error {
	c := &model.UserContestation{
		UserID:      userID,
		Reason:      reason,
		Status:      "pending",
		RequestedAt: time.Now(),
	}
	return s.repo.Create(c)
}

func (s *userContestationService) ListPending(limit, offset int) ([]model.UserContestation, int64, error) {
	return s.repo.ListPending(limit, offset)
}

func (s *userContestationService) Approve(id uint, adminID uint, notes *string) error {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if c.Status != "pending" {
		return nil
	}
	now := time.Now()
	c.Status = "approved"
	c.ReviewedAt = &now
	c.ReviewedBy = &adminID
	c.ReviewNotes = notes
	return s.repo.Update(c)
}

func (s *userContestationService) Reject(id uint, adminID uint, notes *string) error {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if c.Status != "pending" {
		return nil
	}
	now := time.Now()
	c.Status = "rejected"
	c.ReviewedAt = &now
	c.ReviewedBy = &adminID
	c.ReviewNotes = notes
	return s.repo.Update(c)
}
