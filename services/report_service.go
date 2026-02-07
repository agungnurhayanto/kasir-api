package services

import (
	"KASIR-API/dto"
	"KASIR-API/repositories"
	"KASIR-API/utils"
	"time"
)

type ReportService struct {
	repo *repositories.ReportRepository
}

func NewReportService(repo *repositories.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetReport(startDate, endDate string) (*dto.ReportResponse, error) {

	loc := utils.AppLocation

	layout := "2006-01-02"

	startLocal, err := time.ParseInLocation(layout, startDate, loc)
	if err != nil {
		return nil, err
	}

	endLocal, err := time.ParseInLocation(layout, endDate, loc)
	if err != nil {
		return nil, err
	}

	// end date exclusive (besok 00:00)
	endLocal = endLocal.Add(24 * time.Hour)

	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	totalRevenue, totalTransaksi, nama, qty, err :=
		s.repo.GetReport(startUTC, endUTC)

	if err != nil {
		return nil, err
	}

	response := dto.ReportResponse{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: dto.BestProduct{
			Nama:       nama,
			QtyTerjual: qty,
		},
	}

	return &response, nil
}

func (s *ReportService) GetTodayReport() (*dto.ReportResponse, error) {

	loc := utils.AppLocation
	now := time.Now().In(loc)

	startLocal := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0,
		loc,
	)

	endLocal := startLocal.Add(24 * time.Hour)

	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	totalRevenue, totalTransaksi, nama, qty, err :=
		s.repo.GetReport(startUTC, endUTC)

	if err != nil {
		return nil, err
	}

	return &dto.ReportResponse{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransaksi,
		ProdukTerlaris: dto.BestProduct{
			Nama:       nama,
			QtyTerjual: qty,
		},
	}, nil
}

func (s *ReportService) GetTransactionsByRange(startDate, endDate string) (*dto.ReportRangeResponse, error) {

	loc := utils.AppLocation
	layout := "2006-01-02"

	startLocal, err := time.ParseInLocation(layout, startDate, loc)
	if err != nil {
		return nil, err
	}

	endLocal, err := time.ParseInLocation(layout, endDate, loc)
	if err != nil {
		return nil, err
	}

	endLocal = endLocal.Add(24 * time.Hour)

	startUTC := startLocal.UTC()
	endUTC := endLocal.UTC()

	transactions, err := s.repo.GetTransactionsByRange(startUTC, endUTC)
	if err != nil {
		return nil, err
	}

	// convert timezone ke Asia/Jakarta
	for i := range transactions {
		transactions[i].CreatedAt =
			transactions[i].CreatedAt.In(utils.AppLocation)
	}

	return &dto.ReportRangeResponse{
		Transactions: transactions,
	}, nil
}
