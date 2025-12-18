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

// GetMyCart - GET /api/cart
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

// AddToCart - POST /api/cart
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

// UpdateCartItem - PUT /api/cart/:id
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

// RemoveFromCart - DELETE /api/cart/:id
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

// ClearCart - DELETE /api/cart
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

// GetCartItemCount - GET /api/cart/count
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