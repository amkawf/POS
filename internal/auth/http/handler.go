package authhttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"pos-backend/internal/auth/application"
	"pos-backend/internal/pkg/httputil"
)

type Handler struct {
	verifyPinUseCase      *application.VerifyPinUseCase
	listStoreStaffUseCase *application.ListStoreStaffUseCase
	openShiftUseCase      *application.OpenShiftUseCase
	closeShiftUseCase     *application.CloseShiftUseCase
}

func NewHandler(
	verifyPinUseCase *application.VerifyPinUseCase,
	listStoreStaffUseCase *application.ListStoreStaffUseCase,
	openShiftUseCase *application.OpenShiftUseCase,
	closeShiftUseCase *application.CloseShiftUseCase,
) *Handler {
	return &Handler{
		verifyPinUseCase:      verifyPinUseCase,
		listStoreStaffUseCase: listStoreStaffUseCase,
		openShiftUseCase:      openShiftUseCase,
		closeShiftUseCase:     closeShiftUseCase,
	}
}

// RegisterRoutes mendaftarkan endpoint autentikasi & shift ke router Gin
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/auth/pin-login", h.PinLogin)
	rg.GET("/auth/staff", h.ListStoreStaff)
	rg.POST("/auth/shifts/open", h.OpenShift)
	rg.POST("/auth/shifts/close", h.CloseShift)
}

// PinLogin menangani POST /api/v1/auth/pin-login
func (h *Handler) PinLogin(c *gin.Context) {
	var req PinLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	result, err := h.verifyPinUseCase.Execute(c.Request.Context(), req.StoreID, req.PIN)
	if err != nil {
		httputil.Error(c, http.StatusUnauthorized, "AUTH_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, result)
}

// ListStoreStaff menangani GET /api/v1/auth/staff?store_id=...
func (h *Handler) ListStoreStaff(c *gin.Context) {
	storeID, ok := httputil.ParseUUIDQuery(c, "store_id")
	if !ok {
		return
	}

	staffList, err := h.listStoreStaffUseCase.Execute(c.Request.Context(), storeID)
	if err != nil {
		httputil.InternalError(c, "LIST_STAFF_FAILED", err.Error())
		return
	}

	res := make([]StaffProfileResponse, 0, len(staffList))
	for _, s := range staffList {
		res = append(res, StaffProfileResponse{
			ID:   s.ID,
			Name: s.Name,
			Role: s.Role,
		})
	}

	httputil.JSON(c, http.StatusOK, gin.H{"staff": res})
}

// OpenShift menangani POST /api/v1/auth/shifts/open
func (h *Handler) OpenShift(c *gin.Context) {
	var req OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	shift, err := h.openShiftUseCase.Execute(
		c.Request.Context(),
		req.CompanyID,
		req.StoreID,
		req.UserID,
		req.StartingCash,
	)
	if err != nil {
		httputil.InternalError(c, "OPEN_SHIFT_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, shift)
}

// CloseShift menangani POST /api/v1/auth/shifts/close
func (h *Handler) CloseShift(c *gin.Context) {
	var req CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.BadRequest(c, "INVALID_PAYLOAD", err.Error())
		return
	}

	err := h.closeShiftUseCase.Execute(
		c.Request.Context(),
		req.ShiftID,
		req.ActualEndingCash,
		req.ExpectedEndingCash,
		req.Notes,
	)
	if err != nil {
		httputil.InternalError(c, "CLOSE_SHIFT_FAILED", err.Error())
		return
	}

	httputil.JSON(c, http.StatusOK, gin.H{
		"message": "Sesi kasir berhasil ditutup",
	})
}
