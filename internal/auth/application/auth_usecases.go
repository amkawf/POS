package application

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"pos-backend/internal/auth/domain"
	"pos-backend/internal/auth/repository"
	"pos-backend/internal/pkg/jwtutil"
)

// VerifyPinUseCase menangani login instan kasir/koki menggunakan 4-6 digit PIN
type VerifyPinUseCase struct {
	userRepo  repository.UserRepository
	shiftRepo repository.ShiftRepository
	jwtSecret string
}

func NewVerifyPinUseCase(
	userRepo repository.UserRepository,
	shiftRepo repository.ShiftRepository,
	jwtSecret string,
) *VerifyPinUseCase {
	return &VerifyPinUseCase{
		userRepo:  userRepo,
		shiftRepo: shiftRepo,
		jwtSecret: jwtSecret,
	}
}

type PinLoginResult struct {
	User        *domain.User         `json:"user"`
	Token       string               `json:"token"`
	ActiveShift *domain.CashierShift `json:"active_shift,omitempty"`
}

func (uc *VerifyPinUseCase) Execute(
	ctx context.Context,
	storeID uuid.UUID,
	pin string,
) (*PinLoginResult, error) {
	if len(pin) < 4 {
		return nil, errors.New("PIN minimal 4 digit")
	}

	// 1. Cari staf di toko ini yang PIN-nya cocok
	user, err := uc.userRepo.FindByStoreAndPIN(ctx, storeID, pin)
	if err != nil {
		return nil, err
	}

	// 2. Generate token JWT dengan masa aktif 24 jam untuk sesi kasir
	token, err := jwtutil.GenerateToken(user, storeID, uc.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	// 3. Periksa apakah staf kasir ini sudah punya shift laci yang sedang aktif
	var activeShift *domain.CashierShift
	if user.Role == domain.RoleCashier {
		shift, err := uc.shiftRepo.GetActiveShift(ctx, storeID, user.ID)
		if err == nil {
			activeShift = shift
		}
	}

	return &PinLoginResult{
		User:        user,
		Token:       token,
		ActiveShift: activeShift,
	}, nil
}

// ListStoreStaffUseCase mengambil daftar nama staf di suatu toko (untuk tampilan Avatar/Quick Switch)
type ListStoreStaffUseCase struct {
	userRepo repository.UserRepository
}

func NewListStoreStaffUseCase(userRepo repository.UserRepository) *ListStoreStaffUseCase {
	return &ListStoreStaffUseCase{userRepo: userRepo}
}

func (uc *ListStoreStaffUseCase) Execute(ctx context.Context, storeID uuid.UUID) ([]domain.User, error) {
	return uc.userRepo.ListActiveByStore(ctx, storeID)
}

// OpenShiftUseCase membuka sesi laci kasir baru
type OpenShiftUseCase struct {
	shiftRepo repository.ShiftRepository
}

func NewOpenShiftUseCase(shiftRepo repository.ShiftRepository) *OpenShiftUseCase {
	return &OpenShiftUseCase{shiftRepo: shiftRepo}
}

func (uc *OpenShiftUseCase) Execute(
	ctx context.Context,
	companyID, storeID, userID uuid.UUID,
	startingCash float64,
) (*domain.CashierShift, error) {
	// Pastikan kasir belum memiliki shift yang masih terbuka
	existing, err := uc.shiftRepo.GetActiveShift(ctx, storeID, userID)
	if err == nil && existing != nil {
		return existing, nil
	}

	shift := &domain.CashierShift{
		CompanyID:    companyID,
		StoreID:      storeID,
		UserID:       userID,
		StartingCash: startingCash,
		Status:       domain.ShiftStatusOpen,
	}

	if err := uc.shiftRepo.OpenShift(ctx, shift); err != nil {
		return nil, err
	}

	return shift, nil
}

// CloseShiftUseCase menutup sesi kasir dengan perhitungan uang fisik
type CloseShiftUseCase struct {
	shiftRepo repository.ShiftRepository
}

func NewCloseShiftUseCase(shiftRepo repository.ShiftRepository) *CloseShiftUseCase {
	return &CloseShiftUseCase{shiftRepo: shiftRepo}
}

func (uc *CloseShiftUseCase) Execute(
	ctx context.Context,
	shiftID uuid.UUID,
	actualEndingCash float64,
	expectedEndingCash float64,
	notes *string,
) error {
	return uc.shiftRepo.CloseShift(ctx, shiftID, actualEndingCash, expectedEndingCash, notes)
}
