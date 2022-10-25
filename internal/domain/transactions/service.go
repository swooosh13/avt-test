package transactions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
	"github.com/swooosh13/avito-test/internal/repository/s3"
)

type Repository interface {
	Revenue(ctx context.Context, dto *RevenueDTO) error                   //
	Report(ctx context.Context, dto *ReportDTO) ([]TransactReport, error) // отчет в формате csv
	GetTransactionsByUser(ctx context.Context, id int64, params ListParams) ([]GetUserTransactionDTO, int, error)
}

type Service struct {
	repo     Repository
	s3Client *s3.Client
	logger   *zerolog.Logger
}

func NewService(repo Repository, client *s3.Client, logger *zerolog.Logger) *Service {
	return &Service{
		repo:     repo,
		logger:   logger,
		s3Client: client,
	}
}

func (s *Service) GetTransactionsByUser(ctx context.Context, id int64, params ListParams) ([]GetUserTransactionDTO, int, error) {
	transactions, total, err := s.repo.GetTransactionsByUser(ctx, id, params)
	if err != nil {
		return nil, 0, fmt.Errorf("get transactions by user: %w", err)
	}

	return transactions, total, nil
}

func (s *Service) Revenue(ctx context.Context, dto *RevenueDTO) error {
	err := s.repo.Revenue(ctx, dto)
	if err != nil {
		return fmt.Errorf("revenue by params: %w", err)
	}
	return nil
}

func (s *Service) Report(ctx context.Context, dto *ReportDTO) (string, error) {
	reports, err := s.repo.Report(ctx, dto)
	if err != nil {
		return "", fmt.Errorf("get report by params: %w", err)
	}

	fileName, err := CreateReportFile(ctx, reports, dto)
	if err != nil {
		return "", fmt.Errorf("create report file: %w", err)
	}

	absFilename, err := filepath.Abs(fileName)
	if err != nil {
		return "", fmt.Errorf("get abs path: %w", err)
	}

	err = s.s3Client.UploadFile(fileName, absFilename)
	if err != nil {
		return "", fmt.Errorf("upload file to s3: %w", err)
	}

	err = os.Remove(absFilename)
	if err != nil {
		return "", fmt.Errorf("remove file: %w", err)
	}

	link := CreateLink("http", s.s3Client.BucketName, s.s3Client.Endpoint, fileName)

	return link, nil
}
