package reservation

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Repository interface {
	GetReservationByID(ctx context.Context, id int64) (Reservation, error)
	ReserveBalance(ctx context.Context, dto *BalanceReserveDTO) error
}

type Service struct {
	repo   Repository
	logger *zerolog.Logger
}

func NewService(repo Repository, logger *zerolog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) GetReservationByID(ctx context.Context, id int64) (Reservation, error) {
	r, err := s.repo.GetReservationByID(ctx, id)
	if err != nil {
		return Reservation{}, fmt.Errorf("get reservation by id: %w", err)
	}

	return r, nil
}

func (s *Service) ReserveBalance(ctx context.Context, dto *BalanceReserveDTO) error {
	err := s.repo.ReserveBalance(ctx, dto)
	if err != nil {
		return fmt.Errorf("reserve balance: %w", err)
	}

	return nil
}
