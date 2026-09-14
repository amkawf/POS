package domain

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Role merepresentasikan peran staf dalam sistem POS restoran
type Role string

const (
	RoleOwner   Role = "OWNER"   // Pemilik resto: Akses penuh semua toko, omzet, & pengaturan
	RoleManager Role = "MANAGER" // Manajer cabang: Akses inventori, bahan, menu, & operasional
	RoleCashier Role = "CASHIER" // Kasir: Akses buat order, bayar, buka/tutup shift laci kasir
	RoleKitchen Role = "KITCHEN" // Koki: Layar tiket dapur & lapor masak batch
)

// User merepresentasikan entitas akun staf di restoran
type User struct {
	ID           uuid.UUID `json:"id"`
	CompanyID    uuid.UUID `json:"company_id"`
	Name         string    `json:"name"`
	Email        *string   `json:"email,omitempty"`
	PasswordHash *string   `json:"-"` // Disembunyikan dari serialisasi JSON (keamanan)
	PinHash      string    `json:"-"` // Disembunyikan dari serialisasi JSON (keamanan)
	Role         Role      `json:"role"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// IsActive memeriksa apakah akun staf masih aktif bekerja
func (u *User) IsActive() bool {
	return u.Status == "ACTIVE"
}

// CheckPIN memvalidasi kecocokan PIN yang diinput kasir menggunakan bcrypt
func (u *User) CheckPIN(plainPIN string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PinHash), []byte(plainPIN))
	return err == nil
}

// CheckPassword memvalidasi kata sandi akun (biasanya untuk manager/owner)
func (u *User) CheckPassword(plainPassword string) bool {
	if u.PasswordHash == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.PasswordHash), []byte(plainPassword))
	return err == nil
}

// HashPIN membuat hash bcrypt dari PIN 4-6 digit sebelum disimpan ke database
func HashPIN(plainPIN string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPIN), bcrypt.DefaultCost)
	return string(bytes), err
}

// HashPassword membuat hash bcrypt dari kata sandi
func HashPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(bytes), err
}
