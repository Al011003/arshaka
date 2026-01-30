// usecase/availability_usecase.go
package usecase

import (
	res "backend/dto/response/barang"
	"backend/model"
	"backend/repo"
	"errors"
	"fmt"
	"time"
)

type AvailabilityUseCase interface {
	GetCalendarSummary(kodeBarang string, monthStr string) (*res.AvailabilityCalendarResponse, error)
	GetDateDetail(kodeBarang string, dateStr string) (*res.DateDetailResponse, error)
}

type availabilityUseCase struct {
	loanRepo   repo.LoanRepository
	barangRepo repo.BarangRepository
	unitRepo   repo.BarangUnitRepository
}

func NewAvailabilityUseCase(
	loanRepo repo.LoanRepository,
	barangRepo repo.BarangRepository,
	unitRepo repo.BarangUnitRepository,
) AvailabilityUseCase {
	return &availabilityUseCase{
		loanRepo:   loanRepo,
		barangRepo: barangRepo,
		unitRepo:   unitRepo,
	}
}

//
// 🔥 Get Calendar Summary - UNIT BASED
//
func (u *availabilityUseCase) GetCalendarSummary(kodeBarang string, monthStr string) (*res.AvailabilityCalendarResponse, error) {
	// 1️⃣ Get barang by KODE
	barang, err := u.barangRepo.FindByKode(kodeBarang)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	// 2️⃣ Get TOTAL units
	totalUnits, err := u.unitRepo.CountTotal(barang.ID)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Parse month
	monthTime, err := time.Parse("2006-01", monthStr)
	if err != nil {
		monthTime = time.Now()
		monthStr = monthTime.Format("2006-01")
	}

	startDate := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// 4️⃣ Get loan items in date range
	loanItems, err := u.loanRepo.GetLoanItemsInDateRange(barang.ID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Initialize calendar
	calendar := make(map[string]*struct {
		Confirmed int
		Pending   int
		Revision  int
	})

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		calendar[dateStr] = &struct {
			Confirmed int
			Pending   int
			Revision  int
		}{}
	}

	// 6️⃣ Populate calendar (COUNT UNITS yang dipinjam)
	for _, item := range loanItems {
		if item.Loan == nil {
			continue
		}

		// Loop tanggal peminjaman
		for d := item.Loan.StartDate; !d.After(item.Loan.EndDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")

			if data, exists := calendar[dateStr]; exists {
				// ✅ Count units yang assigned
				unitCount := len(item.AssignedUnits)

				// Categorize based on status
				if item.Loan.LoanStatus == "APPROVED" && item.Loan.LoanFlowStatus == "REVISION" {
					data.Revision += unitCount
				} else if item.Loan.LoanStatus == "APPROVED" {
					data.Confirmed += unitCount
				} else if item.Loan.LoanStatus == "PENDING" {
					data.Pending += unitCount
				}
			}
		}
	}

	// 7️⃣ Build response
	days := []res.DayAvailability{}
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		data := calendar[dateStr]

		// Available = Total Units - Confirmed - Revision
		available := int(totalUnits) - data.Confirmed - data.Revision
		if available < 0 {
			available = 0
		}

		status := u.determineStatus(int(totalUnits), data.Confirmed, data.Revision, data.Pending, available)

		days = append(days, res.DayAvailability{
			Date:      dateStr,
			DayOfWeek: d.Format("Monday"),
			Confirmed: data.Confirmed,
			Pending:   data.Pending,
			Revision:  data.Revision,
			Available: available,
			Status:    status,
		})
	}

	return &res.AvailabilityCalendarResponse{
		BarangID:   barang.ID,
		BarangNama: barang.Nama,
		BarangKode: barang.Kode,
		TotalStock: int(totalUnits),
		Month:      monthStr,
		Days:       days,
	}, nil
}

//
// 🔥 Get Date Detail - UNIT BASED
//
func (u *availabilityUseCase) GetDateDetail(kodeBarang string, dateStr string) (*res.DateDetailResponse, error) {
	// 1️⃣ Get barang by KODE
	barang, err := u.barangRepo.FindByKode(kodeBarang)
	if err != nil {
		return nil, errors.New("barang tidak ditemukan")
	}

	// 2️⃣ Get total units
	totalUnits, err := u.unitRepo.CountTotal(barang.ID)
	if err != nil {
		return nil, err
	}

	// 3️⃣ Parse date
	targetDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %v", err)
	}

	// 4️⃣ Get loan items for this date
	loanItems, err := u.loanRepo.GetLoanItemsInDateRange(barang.ID, targetDate, targetDate)
	if err != nil {
		return nil, err
	}

	// 5️⃣ Calculate totals and build loan details
	confirmed := 0
	pending := 0
	revision := 0
	loanMap := make(map[string]res.LoanDetail)

	for _, item := range loanItems {
		if item.Loan == nil {
			continue
		}

		// ✅ Count units yang assigned
		unitCount := len(item.AssignedUnits)

		// Categorize
		if item.Loan.LoanStatus == "APPROVED" && item.Loan.LoanFlowStatus == "REVISION" {
			revision += unitCount
		} else if item.Loan.LoanStatus == "APPROVED" {
			confirmed += unitCount
		} else if item.Loan.LoanStatus == "PENDING" {
			pending += unitCount
		}

		// Add to loan details (avoid duplicates)
		if _, exists := loanMap[item.Loan.LoanCode]; !exists {
			userName := "Unknown"
			if item.Loan.User != nil {
				userName = item.Loan.User.Nama
			}

			isRevised := item.Loan.LoanFlowStatus == "REVISION"
			revisionReason := ""
			if isRevised {
				revisionReason = "Stock tidak mencukupi pada tanggal request awal. Telah disesuaikan menjadi " +
					item.Loan.StartDate.Format("02 Jan") + "-" + item.Loan.EndDate.Format("02 Jan")
			}

			priority := u.calculatePriority(item.Loan)

			// ✅ Build unit details
			unitDetails := []res.UnitDetail{}
			for _, unit := range item.AssignedUnits {
				unitDetails = append(unitDetails, res.UnitDetail{
					UnitID:   unit.ID,
					KodeUnit: unit.KodeUnit,
					Kondisi:  unit.Kondisi,
					Status:   unit.Status,
				})
			}

			loanMap[item.Loan.LoanCode] = res.LoanDetail{
				LoanCode:       item.Loan.LoanCode,
				UserName:       userName,
				Quantity:       unitCount,
				Units:          unitDetails,
				LoanStatus:     item.Loan.LoanStatus,
				LoanFlowStatus: item.Loan.LoanFlowStatus,
				StartDate:      item.Loan.StartDate.Format("2006-01-02"),
				EndDate:        item.Loan.EndDate.Format("2006-01-02"),
				IsRevised:      isRevised,
				RevisionReason: revisionReason,
				Reason:         item.Loan.Reason,
				CreatedAt:      item.Loan.CreatedAt.Format("2006-01-02 15:04"),
				Priority:       priority,
			}
		}
	}

	// 6️⃣ Convert map to slice
	loans := make([]res.LoanDetail, 0, len(loanMap))
	for _, loan := range loanMap {
		loans = append(loans, loan)
	}

	// 7️⃣ Calculate available
	available := int(totalUnits) - confirmed - revision
	if available < 0 {
		available = 0
	}

	// 8️⃣ Build summary
	revisionCount := 0
	pendingCount := 0
	for _, loan := range loans {
		if loan.LoanFlowStatus == "REVISION" {
			revisionCount++
		}
		if loan.LoanStatus == "PENDING" {
			pendingCount++
		}
	}

	canAccommodateAll := (confirmed + revision + pending) <= int(totalUnits)
	notes := u.generateNotes(revisionCount, pendingCount, confirmed+revision+pending, int(totalUnits))

	return &res.DateDetailResponse{
		BarangID:   barang.ID,
		BarangNama: barang.Nama,
		BarangKode: barang.Kode,
		TotalStock: int(totalUnits),
		Date:       dateStr,
		Confirmed:  confirmed,
		Pending:    pending,
		Revision:   revision,
		Available:  available,
		Loans:      loans,
		Summary: res.DetailSummary{
			TotalRequests:     len(loans),
			RevisionCount:     revisionCount,
			PendingCount:      pendingCount,
			CanAccommodateAll: canAccommodateAll,
			Notes:             notes,
		},
	}, nil
}

//
// 🔥 Helper Functions
//
func (u *availabilityUseCase) determineStatus(total, confirmed, revision, pending, available int) string {
	if confirmed+revision >= total {
		return "fully_booked"
	}
	if revision > 0 {
		return "revision_pending"
	}
	if pending > 0 && pending >= available {
		return "high_demand"
	}
	if available <= total/4 {
		return "limited"
	}
	return "available"
}

func (u *availabilityUseCase) calculatePriority(loan *model.Loan) string {
	if loan.LoanStatus == "APPROVED" && loan.LoanFlowStatus == "REVISION" {
		return "HIGH"
	}
	return "NORMAL"
}

func (u *availabilityUseCase) generateNotes(revisionCount, pendingCount, totalRequest, totalStock int) string {
	notes := ""

	if revisionCount > 0 {
		notes += fmt.Sprintf("Tanggal ini memiliki %d peminjaman dalam proses revisi (prioritas tinggi)", revisionCount)
	}

	if pendingCount > 0 {
		if notes != "" {
			notes += " dan "
		}
		notes += fmt.Sprintf("%d request pending", pendingCount)
	}

	if totalRequest > totalStock {
		if notes != "" {
			notes += ". "
		}
		notes += "Total request melebihi kapasitas unit."
	}

	if notes == "" {
		notes = "Tidak ada masalah ketersediaan pada tanggal ini."
	}

	return notes
}