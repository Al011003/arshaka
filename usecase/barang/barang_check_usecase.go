package usecase

import (
	res "backend/dto/response/barang"
	"backend/repo"
	"fmt"
	"time"
)

type AvailabilityUseCase interface {
	GetCalendarSummary(barangID uint, monthStr string) (*res.AvailabilityCalendarResponse, error)
}

type availabilityUseCase struct {
	loanRepo repo.LoanRepository
	barangRepo repo.BarangRepository
}

func NewAvailabilityUseCase(loanrRepo repo.LoanRepository, barangRepo repo.BarangRepository) AvailabilityUseCase {
	return &availabilityUseCase{loanRepo: loanrRepo, barangRepo : barangRepo}
}

func (u *availabilityUseCase) GetCalendarSummary(barangID uint, monthStr string) (*res.AvailabilityCalendarResponse, error) {
	// Get barang
	barang, err := u.barangRepo.FindByID(barangID)
	if err != nil {
		return nil, err
	}

	fmt.Println("BARANG ID:", barangID)
	// Parse month
	monthTime, err := time.Parse("2006-01", monthStr)
	if err != nil {
		monthTime = time.Now()
		monthStr = monthTime.Format("2006-01")
	}

	// Get first and last day of month
	startDate := time.Date(monthTime.Year(), monthTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	// Get loan items
	loanItems, err := u.loanRepo.GetLoanItemsInDateRange(barangID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Initialize calendar
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

	// Populate calendar
	for _, item := range loanItems {
		if item.Loan == nil {
			continue
		}

		for d := item.Loan.StartDate; !d.After(item.Loan.EndDate); d = d.AddDate(0, 0, 1) {
			dateStr := d.Format("2006-01-02")

			if data, exists := calendar[dateStr]; exists {
				// Categorize based on status
				if item.Loan.LoanStatus == "APPROVED" && item.Loan.LoanFlowStatus == "REVISION" {
					// APPROVED tapi masih REVISION (waiting preparation)
					data.Revision += item.Quantity
				} else if item.Loan.LoanStatus == "APPROVED" {
					// APPROVED (REQUESTED, READY, TAKEN)
					data.Confirmed += item.Quantity
				} else if item.Loan.LoanStatus == "PENDING" {
					// PENDING
					data.Pending += item.Quantity
				}
			}
		}
	}

	// Build response
	days := []res.DayAvailability{}
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		data := calendar[dateStr]

		// Available = Total - Confirmed - Revision
		available := barang.StokTotal - data.Confirmed - data.Revision
		if available < 0 {
			available = 0
		}

		// Determine status
		status := u.determineStatus(barang.StokTotal, data.Confirmed, data.Revision, data.Pending, available)

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
		TotalStock: barang.StokTotal,
		Month:      monthStr,
		Days:       days,
	}, nil
}



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

// func (u *availabilityUseCase) calculatePriority(loan *model.Loan) string {
// 	if loan.LoanStatus == "APPROVED" && loan.LoanFlowStatus == "REVISION" {
// 		return "HIGH"
// 	}
// 	return "NORMAL"
// }

// func (u *availabilityUseCase) generateNotes(revisionCount, pendingCount, totalRequest, totalStock int) string {
// 	notes := ""
	
// 	if revisionCount > 0 {
// 		notes += fmt.Sprintf("Tanggal ini memiliki %d peminjaman dalam proses revisi (prioritas tinggi)", revisionCount)
// 	}
	
// 	if pendingCount > 0 {
// 		if notes != "" {
// 			notes += " dan "
// 		}
// 		notes += fmt.Sprintf("%d request pending", pendingCount)
// 	}
	
// 	if totalRequest > totalStock {
// 		if notes != "" {
// 			notes += ". "
// 		}
// 		notes += "Total request melebihi kapasitas stok."
// 	}
	
// 	if notes == "" {
// 		notes = "Tidak ada masalah ketersediaan pada tanggal ini."
// 	}
	
// 	return notes
// }

// func (u *availabilityUseCase) GetDateDetail(barangID uint, dateStr string) (*res.DateDetailResponse, error) {
// 	// Get barang
// 	barang, err := u.barangRepo.FindByID(barangID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Parse date
// 	targetDate, err := time.Parse("2006-01-02", dateStr)
// 	if err != nil {
// 		return nil, fmt.Errorf("invalid date format: %v", err)
// 	}

// 	// Get loan items for this date
// 	loanItems, err := u.loanRepo.GetLoanItemsInDateRange(barangID, targetDate, targetDate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Calculate totals and build loan details
// 	confirmed := 0
// 	pending := 0
// 	revision := 0
// 	loanMap := make(map[string]res.LoanDetail)

// 	for _, item := range loanItems {
// 		if item.Loan == nil {
// 			continue
// 		}

// 		// Count quantities
// 		if item.Loan.LoanStatus == "APPROVED" && item.Loan.LoanFlowStatus == "REVISION" {
// 			revision += item.Quantity
// 		} else if item.Loan.LoanStatus == "APPROVED" {
// 			confirmed += item.Quantity
// 		} else if item.Loan.LoanStatus == "PENDING" {
// 			pending += item.Quantity
// 		}

// 		// Add to loan details (avoid duplicates)
// 		if _, exists := loanMap[item.Loan.LoanCode]; !exists {
// 			userName := "Unknown"
// 			if item.Loan.User != nil {
// 				userName = item.Loan.User.Name
// 			}

// 			isRevised := item.Loan.LoanFlowStatus == "REVISION"
// 			revisionReason := ""
// 			if isRevised {
// 				revisionReason = "Stock tidak mencukupi pada tanggal request awal. Telah disesuaikan menjadi " + 
// 					item.Loan.StartDate.Format("02 Jan") + "-" + item.Loan.EndDate.Format("02 Jan")
// 			}

// 			priority := u.calculatePriority(item.Loan)

// 			loanMap[item.Loan.LoanCode] = res.LoanDetail{
// 				LoanCode:       item.Loan.LoanCode,
// 				UserName:       userName,
// 				Quantity:       item.Quantity,
// 				LoanStatus:     item.Loan.LoanStatus,
// 				LoanFlowStatus: item.Loan.LoanFlowStatus,
// 				StartDate:      item.Loan.StartDate.Format("2006-01-02"),
// 				EndDate:        item.Loan.EndDate.Format("2006-01-02"),
// 				IsRevised:      isRevised,
// 				RevisionReason: revisionReason,
// 				Reason:         item.Loan.Reason,
// 				CreatedAt:      item.Loan.CreatedAt.Format("2006-01-02 15:04"),
// 				Priority:       priority,
// 			}
// 		}
// 	}

// 	// Convert map to slice
// 	loans := make([]res.LoanDetail, 0, len(loanMap))
// 	for _, loan := range loanMap {
// 		loans = append(loans, loan)
// 	}

// 	// Sort: HIGH priority first, then NORMAL
// 	sort.Slice(loans, func(i, j int) bool {
// 		if loans[i].Priority != loans[j].Priority {
// 			return loans[i].Priority == "HIGH"
// 		}
// 		return loans[i].CreatedAt < loans[j].CreatedAt
// 	})

// 	available := barang.StokTotal - confirmed - revision
// 	if available < 0 {
// 		available = 0
// 	}

// 	// Build summary
// 	revisionCount := 0
// 	pendingCount := 0
// 	for _, loan := range loans {
// 		if loan.LoanFlowStatus == "REVISION" {
// 			revisionCount++
// 		}
// 		if loan.LoanStatus == "PENDING" {
// 			pendingCount++
// 		}
// 	}

// 	canAccommodateAll := (confirmed + revision + pending) <= barang.StokTotal
// 	notes := u.generateNotes(revisionCount, pendingCount, confirmed+revision+pending, barang.StokTotal)

// 	return &res.DateDetailResponse{
// 		BarangID:   barang.ID,
// 		BarangNama: barang.Nama,
// 		BarangKode: barang.Kode,
// 		TotalStock: barang.StokTotal,
// 		Date:       dateStr,
// 		Confirmed:  confirmed,
// 		Pending:    pending,
// 		Revision:   revision,
// 		Available:  available,
// 		Loans:      loans,
// 		Summary: res.DetailSummary{
// 			TotalRequests:     len(loans),
// 			RevisionCount:     revisionCount,
// 			PendingCount:      pendingCount,
// 			CanAccommodateAll: canAccommodateAll,
// 			Notes:             notes,
// 		},
// 	}, nil
// }