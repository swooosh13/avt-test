package transactions

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func CreateReportFile(ctx context.Context, reports []TransactReport, dto *ReportDTO) (string, error) {
	today := time.Now()
	fileName := fmt.Sprintf("report_%d%d_%d%d%d.csv", dto.Period.Year(), dto.Period.Month(), today.Year(), today.Month(), today.Day())

	f, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}

	wr := csv.NewWriter(f)
	err = wr.Write([]string{"service_id", "total amount"})
	if err != nil {
		return "", fmt.Errorf("write csv header: %w", err)
	}

	for _, v := range reports {
		svcID := strconv.Itoa(int(v.ServiceID))
		totalAmount := strconv.Itoa(int(v.TotalAmount))
		err := wr.Write([]string{svcID, totalAmount})
		if err != nil {
			return "", fmt.Errorf("write to csv: %w", err)
		}
	}
	wr.Flush()
	f.Close()

	return fileName, nil
}

func CreateLink(protocol, bucket, endpoint, fileName string) string {
	return fmt.Sprintf("%s://%s.%s/%s", protocol, bucket, endpoint, fileName)
}
