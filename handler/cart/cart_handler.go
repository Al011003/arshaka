// handler/cart_handler.go
package handler

import (
	"strconv"

	cart "backend/dto/request/cart"
	usecase "backend/usecase/user/cart"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	cartUsecase usecase.CartUsecase
}

func NewCartHandler(cartUC usecase.CartUsecase) *CartHandler {
	return &CartHandler{
		cartUsecase: cartUC,
	}
}

// GetMyCart godoc
// @Summary      Get my cart
// @Description  Mengambil semua item cart milik user yang sedang login
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Berhasil mengambil cart"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/cart [get]
// @Security     ApiKeyAuth
func (h *CartHandler) GetMyCart(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Get cart
	result, err := h.cartUsecase.GetMyCart(userID.(uint))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "berhasil mengambil cart")
}

// AddToCart godoc
// @Summary      Add item to cart
// @Description  Menambahkan barang ke dalam cart user
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        request  body      cart.AddToCartRequest  true  "Data cart"
// @Success      201      {object}  map[string]interface{}  "Berhasil menambahkan ke cart"
// @Failure      400      {object}  map[string]interface{}  "Bad request - validation error"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /api/cart [post]
// @Security     ApiKeyAuth
func (h *CartHandler) AddToCart(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Validate request
	var req cart.AddToCartRequest
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Add to cart
	result, err := h.cartUsecase.AddToCart(userID.(uint), req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, result, "berhasil menambahkan ke cart")
}

// UpdateCartItem godoc
// @Summary      Update cart item
// @Description  Mengupdate quantity atau data item di cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id       path      string                     true  "Cart Item ID"
// @Param        request  body      cart.UpdateCartItemRequest true  "Update cart item"
// @Success      200      {object}  map[string]interface{}  "Berhasil update cart item"
// @Failure      400      {object}  map[string]interface{}  "Bad request"
// @Failure      401      {object}  map[string]interface{}  "Unauthorized"
// @Failure      500      {object}  map[string]interface{}  "Internal server error"
// @Router       /api/cart/{id} [put]
// @Security     ApiKeyAuth
func (h *CartHandler) UpdateCartItem(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Get cart item ID from URL param
	id := c.Param("id")
	cartItemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID cart item tidak valid")
		return
	}

	// Validate request
	var req cart.UpdateCartItemRequest
	if err := req.BindAndValidate(c); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// Update cart item
	result, err := h.cartUsecase.UpdateCartItem(userID.(uint), uint(cartItemID), req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "berhasil update cart item")
}

// RemoveFromCart godoc
// @Summary      Remove item from cart
// @Description  Menghapus satu item dari cart user
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Cart Item ID"
// @Success      200  {object}  map[string]interface{}  "Berhasil menghapus item"
// @Failure      400  {object}  map[string]interface{}  "Bad request"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/cart/{id} [delete]
// @Security     ApiKeyAuth
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Get cart item ID from URL param
	id := c.Param("id")
	cartItemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.BadRequest(c, "ID cart item tidak valid")
		return
	}

	// Remove from cart
	err = h.cartUsecase.RemoveFromCart(userID.(uint), uint(cartItemID))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "berhasil menghapus dari cart")
}

// ClearCart godoc
// @Summary      Clear cart
// @Description  Menghapus seluruh item cart milik user
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Berhasil mengosongkan cart"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/cart [delete]
// @Security     ApiKeyAuth
func (h *CartHandler) ClearCart(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Clear cart
	err := h.cartUsecase.ClearMyCart(userID.(uint))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, nil, "berhasil mengosongkan cart")
}

// GetCartItemCount godoc
// @Summary      Get cart item count
// @Description  Mengambil jumlah total item di cart user
// @Tags         cart
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "Berhasil mengambil jumlah item"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Failure      500  {object}  map[string]interface{}  "Internal server error"
// @Router       /api/cart/count [get]
// @Security     ApiKeyAuth
func (h *CartHandler) GetCartItemCount(c *gin.Context) {
	// Get user ID from JWT middleware context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "user tidak terautentikasi")
		return
	}

	// Get cart item count
	result, err := h.cartUsecase.GetCartItemCount(userID.(uint))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, result, "berhasil mengambil jumlah item")
}