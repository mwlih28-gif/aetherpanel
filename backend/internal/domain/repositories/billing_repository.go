package repositories

import (
	"context"
	"time"

	"github.com/aetherpanel/aether-panel/internal/domain/entities"
	"github.com/google/uuid"
)

// TransactionRepository defines the interface for transaction data access
type TransactionRepository interface {
	Create(ctx context.Context, tx *entities.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Transaction, error)
	Update(ctx context.Context, tx *entities.Transaction) error
	GetByUserID(ctx context.Context, userID uuid.UUID, params ListParams) ([]*entities.Transaction, int64, error)
	GetByReference(ctx context.Context, reference string) (*entities.Transaction, error)
	GetByDateRange(ctx context.Context, start, end time.Time) ([]*entities.Transaction, error)
	SumByUserID(ctx context.Context, userID uuid.UUID, txType entities.TransactionType) (float64, error)
}

// PackageRepository defines the interface for package data access
type PackageRepository interface {
	Create(ctx context.Context, pkg *entities.Package) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Package, error)
	GetBySlug(ctx context.Context, slug string) (*entities.Package, error)
	Update(ctx context.Context, pkg *entities.Package) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) ([]*entities.Package, int64, error)
	GetActive(ctx context.Context) ([]*entities.Package, error)
	GetByCategory(ctx context.Context, category string) ([]*entities.Package, error)
}

// SubscriptionRepository defines the interface for subscription data access
type SubscriptionRepository interface {
	Create(ctx context.Context, sub *entities.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Subscription, error)
	Update(ctx context.Context, sub *entities.Subscription) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*entities.Subscription, error)
	GetByServerID(ctx context.Context, serverID uuid.UUID) (*entities.Subscription, error)
	GetActive(ctx context.Context) ([]*entities.Subscription, error)
	GetExpiring(ctx context.Context, before time.Time) ([]*entities.Subscription, error)
	GetExpired(ctx context.Context) ([]*entities.Subscription, error)
	Cancel(ctx context.Context, id uuid.UUID) error
	Suspend(ctx context.Context, id uuid.UUID) error
	Renew(ctx context.Context, id uuid.UUID, newEndDate time.Time) error
}

// InvoiceRepository defines the interface for invoice data access
type InvoiceRepository interface {
	Create(ctx context.Context, invoice *entities.Invoice) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Invoice, error)
	GetByNumber(ctx context.Context, number string) (*entities.Invoice, error)
	Update(ctx context.Context, invoice *entities.Invoice) error
	GetByUserID(ctx context.Context, userID uuid.UUID, params ListParams) ([]*entities.Invoice, int64, error)
	GetOverdue(ctx context.Context) ([]*entities.Invoice, error)
	MarkPaid(ctx context.Context, id uuid.UUID) error
	GenerateNumber(ctx context.Context) (string, error)
}

// CouponRepository defines the interface for coupon data access
type CouponRepository interface {
	Create(ctx context.Context, coupon *entities.Coupon) error
	GetByID(ctx context.Context, id uuid.UUID) (*entities.Coupon, error)
	GetByCode(ctx context.Context, code string) (*entities.Coupon, error)
	Update(ctx context.Context, coupon *entities.Coupon) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListParams) ([]*entities.Coupon, int64, error)
	IncrementUsage(ctx context.Context, id uuid.UUID) error
	GetValidForPackage(ctx context.Context, packageID uuid.UUID) ([]*entities.Coupon, error)
}
