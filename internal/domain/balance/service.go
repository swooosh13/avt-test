package balance

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

type Repository interface {
	GetBalanceByID(ctx context.Context, id int64) (int64, error)
	UpdateBalance(ctx context.Context, dto *UpdateBalanceDTO) error
	TransferFromTo(ctx context.Context, dto *TransferFromToDTO) error
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

func (s *Service) GetBalanceByID(ctx context.Context, id int64) (int64, error) {
	b, err := s.repo.GetBalanceByID(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("get balance by id: %w", err)
	}

	return b, nil
}

func (s *Service) UpdateBalance(ctx context.Context, dto *UpdateBalanceDTO) error {
	err := s.repo.UpdateBalance(ctx, dto)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}

func (s *Service) TransferFromTo(ctx context.Context, dto *TransferFromToDTO) error {
	err := s.repo.TransferFromTo(ctx, dto)
	if err != nil {
		return fmt.Errorf("transfer from to: %w", err)
	}

	return nil
}
