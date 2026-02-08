package handler

import (
	req "backend/dto/request/loan"
	usecase "backend/usecase/loan"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type LoanUserHandler struct {
	loanUsecase usecase.LoanUserUsecase
}

func NewLoanUserHandler(loanUC usecase.LoanUserUsecase) *LoanUserHandler {
	return &LoanUserHandler{
		loanUsecase: loanUC,
	}
}

// CreateLoan godoc
// @Summary      Create loan
// @Description  User membuat pengajuan peminjaman barang. Sistem akan otomatis mengecek availability dengan gap 1 hari rule dan memberikan saran alternatif jika barang tidak tersedia.
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        request  body      request.CreateLoanRequest  true  "Create loan request"
// @Success      201      {object}  map[string]interface{}{data=response.LoanResponse}  "Loan berhasil dibuat"
// @Failure      400      {object}  map[string]interface{}{error=string,suggestions=[]string}  "Bad request - Barang tidak tersedia (dengan saran alternatif)"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/user/loan [post]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	var reqBody req.CreateLoanRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.loanUsecase.Create(reqBody, userID.(uint))
	if err != nil {
		// Error message udah include suggestions dari usecase
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Created(c, result, "peminjaman berhasil dibuat")
}

// UpdateLoanHeader godoc
// @Summary      Update loan header
// @Description  Update data header peminjaman (tanggal, catatan). Hanya bisa dilakukan jika loan status = REQUESTED atau REVISION. Sistem akan mengecek availability dengan gap 1 hari rule untuk tanggal baru.
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        loan_code  path      string                      true  "Loan code"
// @Param        request    body      request.UpdateLoanRequest   true  "Update loan header"
// @Success      200        {object}  map[string]interface{}  "Peminjaman berhasil diperbarui"
// @Failure      400        {object}  map[string]interface{}{error=string}  "Bad request - Loan tidak bisa diedit atau barang tidak tersedia di tanggal baru"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      403        {object}  map[string]interface{}  "Forbidden - Bukan pemilik loan"
// @Router       /api/user/loan/{loan_code} [put]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) UpdateHeader(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	var reqBody req.UpdateLoanRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.loanUsecase.UpdateHeader(loanCode, userID, reqBody); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "peminjaman berhasil diperbarui")
}

// UpdateLoanItems godoc
// @Summary      Update loan items
// @Description  Update daftar barang dalam peminjaman (ADD/UPDATE/DELETE). Hanya bisa dilakukan jika loan status = REQUESTED atau REVISION. Sistem akan mengecek availability dengan gap 1 hari rule untuk setiap perubahan.
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        loan_code  path      string                            true  "Loan code"
// @Param        request    body      request.UpdateLoanItemsRequest    true  "Update loan items (action: ADD/UPDATE/DELETE)"
// @Success      200        {object}  map[string]interface{}  "Barang peminjaman berhasil diperbarui"
// @Failure      400        {object}  map[string]interface{}{error=string}  "Bad request - Loan tidak bisa diedit atau barang tidak tersedia"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      403        {object}  map[string]interface{}  "Forbidden - Bukan pemilik loan"
// @Router       /api/user/loan/{loan_code}/items [put]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) UpdateItems(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	var reqBody req.UpdateLoanItemsRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.loanUsecase.UpdateItems(loanCode, userID, reqBody); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "barang peminjaman berhasil diperbarui")
}

// GetMyLoans godoc
// @Summary      Get my loans
// @Description  Mengambil daftar semua peminjaman milik user yang sedang login
// @Tags         user-loan
// @Produce      json
// @Success      200  {object}  map[string]interface{}{data=[]response.LoanListResponse}  "Daftar peminjaman"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/user/loan [get]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) GetMyLoans(c *gin.Context) {
	userID := c.GetUint("user_id")

	loans, err := h.loanUsecase.GetMyLoans(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, loans, "daftar peminjaman")
}

// CancelLoan godoc
// @Summary      Cancel loan
// @Description  Membatalkan peminjaman (hanya bisa dilakukan jika loan_status = PENDING). Loan akan di-soft delete dan status berubah menjadi REJECTED.
// @Tags         user-loan
// @Produce      json
// @Param        loan_code  path      string  true  "Loan code"
// @Success      200        {object}  map[string]interface{}  "Peminjaman berhasil dibatalkan"
// @Failure      400        {object}  map[string]interface{}{error=string}  "Bad request - Loan tidak bisa dibatalkan"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      403        {object}  map[string]interface{}  "Forbidden - Bukan pemilik loan"
// @Router       /api/user/loan/{loan_code} [delete]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	if err := h.loanUsecase.Cancel(loanCode, userID); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, nil, "peminjaman berhasil dibatalkan")
}

// GetLoanDetail godoc
// @Summary      Get loan detail
// @Description  Mengambil detail lengkap peminjaman berdasarkan loan code, termasuk daftar barang dan unit yang dipinjam
// @Tags         user-loan
// @Produce      json
// @Param        loan_code  path      string  true  "Loan code"
// @Success      200        {object}  map[string]interface{}{data=response.LoanDetailResponse}  "Detail peminjaman"
// @Failure      400        {object}  map[string]interface{}  "Bad request - Loan code tidak valid"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      403        {object}  map[string]interface{}  "Forbidden - Bukan pemilik loan"
// @Failure      404        {object}  map[string]interface{}  "Not found - Loan tidak ditemukan"
// @Router       /api/user/loan/{loan_code} [get]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")

	loanCode := c.Param("loan_code")
	if loanCode == "" {
		utils.BadRequest(c, "loan_code tidak valid")
		return
	}

	result, err := h.loanUsecase.GetDetail(loanCode, userID)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, result, "detail peminjaman")
}

// 🔥 NEW: CheckAvailability godoc
// @Summary      Check barang availability
// @Description  Mengecek ketersediaan barang untuk tanggal tertentu dengan gap 1 hari rule. Jika tidak tersedia, akan memberikan saran tanggal alternatif dan rekomendasi barang serupa.
// @Tags         user-loan
// @Accept       json
// @Produce      json
// @Param        request  body      request.CheckAvailabilityRequest  true  "Check availability request"
// @Success      200      {object}  map[string]interface{}{data=response.AvailabilityCheckResponse}  "Hasil pengecekan availability"
// @Failure      400      {object}  map[string]interface{}  "Bad request"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Router       /api/user/loan/check-availability [post]
// @Security     ApiKeyAuth
func (h *LoanUserHandler) CheckAvailability(c *gin.Context) {
	var reqBody req.CheckAvailabilityRequest
	if err := reqBody.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	result, err := h.loanUsecase.CheckAvailabilityWithSuggestions(
		reqBody.KodeBarang,
		reqBody.Quantity,
		reqBody.ParsedStartDate,
		reqBody.ParsedEndDate,
	)

	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	utils.Success(c, result, "hasil pengecekan ketersediaan")
}