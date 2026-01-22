package entities

import (
	"time"

	"github.com/google/uuid"
)

// TransactionType represents the type of billing transaction
type TransactionType string

const (
	TransactionTypeCredit   TransactionType = "credit"
	TransactionTypeDebit    TransactionType = "debit"
	TransactionTypeRefund   TransactionType = "refund"
	TransactionTypeTransfer TransactionType = "transfer"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

// PaymentMethod represents a payment method type
type PaymentMethod string

const (
	PaymentMethodStripe    PaymentMethod = "stripe"
	PaymentMethodPayPal    PaymentMethod = "paypal"
	PaymentMethodCrypto    PaymentMethod = "crypto"
	PaymentMethodManual    PaymentMethod = "manual"
	PaymentMethodInternal  PaymentMethod = "internal"
)

// Transaction represents a billing transaction
type Transaction struct {
	ID              uuid.UUID         `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID         `json:"user_id" gorm:"type:uuid;not null;index"`
	User            *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type            TransactionType   `json:"type" gorm:"type:varchar(20);not null"`
	Status          TransactionStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	Amount          float64           `json:"amount" gorm:"type:decimal(12,2);not null"`
	Currency        string            `json:"currency" gorm:"size:3;default:'USD'"`
	Description     string            `json:"description" gorm:"size:500"`
	Reference       string            `json:"reference" gorm:"size:100;index"` // External reference
	PaymentMethod   PaymentMethod     `json:"payment_method" gorm:"type:varchar(20)"`
	PaymentDetails  map[string]interface{} `json:"payment_details" gorm:"type:jsonb;default:'{}'"`
	BalanceBefore   float64           `json:"balance_before" gorm:"type:decimal(12,2)"`
	BalanceAfter    float64           `json:"balance_after" gorm:"type:decimal(12,2)"`
	ProcessedAt     *time.Time        `json:"processed_at"`
	ProcessedBy     *uuid.UUID        `json:"processed_by" gorm:"type:uuid"`
	Notes           string            `json:"notes" gorm:"size:500"`
	CreatedAt       time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Transaction
func (Transaction) TableName() string {
	return "transactions"
}

// Package represents a service package/plan
type Package struct {
	ID              uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name            string    `json:"name" gorm:"not null;size:100"`
	Slug            string    `json:"slug" gorm:"uniqueIndex;not null;size:50"`
	Description     string    `json:"description" gorm:"type:text"`
	Category        string    `json:"category" gorm:"size:50;index"`
	
	// Pricing
	PriceMonthly    float64   `json:"price_monthly" gorm:"type:decimal(10,2);default:0"`
	PriceQuarterly  float64   `json:"price_quarterly" gorm:"type:decimal(10,2);default:0"`
	PriceYearly     float64   `json:"price_yearly" gorm:"type:decimal(10,2);default:0"`
	SetupFee        float64   `json:"setup_fee" gorm:"type:decimal(10,2);default:0"`
	Currency        string    `json:"currency" gorm:"size:3;default:'USD'"`

	// Resources
	MemoryLimit     int64     `json:"memory_limit" gorm:"default:1024"`     // MB
	DiskLimit       int64     `json:"disk_limit" gorm:"default:10240"`      // MB
	CPULimit        int       `json:"cpu_limit" gorm:"default:100"`         // Percentage
	DatabaseLimit   int       `json:"database_limit" gorm:"default:1"`
	AllocationLimit int       `json:"allocation_limit" gorm:"default:1"`
	BackupLimit     int       `json:"backup_limit" gorm:"default:2"`
	ServerLimit     int       `json:"server_limit" gorm:"default:1"`

	// Game restrictions
	AllowedGames    []uuid.UUID `json:"allowed_games" gorm:"type:jsonb"`
	AllowedEggs     []uuid.UUID `json:"allowed_eggs" gorm:"type:jsonb"`
	AllowedNodes    []uuid.UUID `json:"allowed_nodes" gorm:"type:jsonb"`

	// Display
	Features        []string  `json:"features" gorm:"type:jsonb"`
	Badge           string    `json:"badge" gorm:"size:50"` // e.g., "Popular", "Best Value"
	Color           string    `json:"color" gorm:"size:7"`
	SortOrder       int       `json:"sort_order" gorm:"default:0"`
	IsActive        bool      `json:"is_active" gorm:"default:true"`
	IsHidden        bool      `json:"is_hidden" gorm:"default:false"`
	
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`
}

// TableName returns the table name for Package
func (Package) TableName() string {
	return "packages"
}

// Subscription represents a user's subscription to a package
type Subscription struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PackageID       uuid.UUID  `json:"package_id" gorm:"type:uuid;not null;index"`
	Package         *Package   `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	ServerID        *uuid.UUID `json:"server_id" gorm:"type:uuid;index"`
	Server          *Server    `json:"server,omitempty" gorm:"foreignKey:ServerID"`
	
	// Billing
	BillingCycle    string     `json:"billing_cycle" gorm:"size:20"` // monthly, quarterly, yearly
	Amount          float64    `json:"amount" gorm:"type:decimal(10,2)"`
	Currency        string     `json:"currency" gorm:"size:3;default:'USD'"`
	
	// Status
	Status          string     `json:"status" gorm:"size:20;default:'active'"` // active, suspended, cancelled, expired
	AutoRenew       bool       `json:"auto_renew" gorm:"default:true"`
	
	// Dates
	StartDate       time.Time  `json:"start_date"`
	EndDate         time.Time  `json:"end_date"`
	NextBillingDate *time.Time `json:"next_billing_date"`
	CancelledAt     *time.Time `json:"cancelled_at"`
	SuspendedAt     *time.Time `json:"suspended_at"`
	
	// Grace period
	GracePeriodDays int        `json:"grace_period_days" gorm:"default:3"`
	
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Subscription
func (Subscription) TableName() string {
	return "subscriptions"
}

// IsExpired checks if the subscription is expired
func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.EndDate)
}

// Invoice represents a billing invoice
type Invoice struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InvoiceNumber   string     `json:"invoice_number" gorm:"uniqueIndex;not null;size:50"`
	UserID          uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	User            *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	SubscriptionID  *uuid.UUID `json:"subscription_id" gorm:"type:uuid"`
	Subscription    *Subscription `json:"subscription,omitempty" gorm:"foreignKey:SubscriptionID"`
	
	// Amounts
	Subtotal        float64    `json:"subtotal" gorm:"type:decimal(10,2)"`
	Tax             float64    `json:"tax" gorm:"type:decimal(10,2);default:0"`
	Discount        float64    `json:"discount" gorm:"type:decimal(10,2);default:0"`
	Total           float64    `json:"total" gorm:"type:decimal(10,2)"`
	Currency        string     `json:"currency" gorm:"size:3;default:'USD'"`
	
	// Status
	Status          string     `json:"status" gorm:"size:20;default:'pending'"` // pending, paid, overdue, cancelled
	
	// Dates
	IssueDate       time.Time  `json:"issue_date"`
	DueDate         time.Time  `json:"due_date"`
	PaidAt          *time.Time `json:"paid_at"`
	
	// Items
	Items           []InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID"`
	
	// Notes
	Notes           string     `json:"notes" gorm:"type:text"`
	
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Invoice
func (Invoice) TableName() string {
	return "invoices"
}

// InvoiceItem represents a line item in an invoice
type InvoiceItem struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	InvoiceID   uuid.UUID `json:"invoice_id" gorm:"type:uuid;not null;index"`
	Description string    `json:"description" gorm:"not null;size:500"`
	Quantity    int       `json:"quantity" gorm:"default:1"`
	UnitPrice   float64   `json:"unit_price" gorm:"type:decimal(10,2)"`
	Total       float64   `json:"total" gorm:"type:decimal(10,2)"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// TableName returns the table name for InvoiceItem
func (InvoiceItem) TableName() string {
	return "invoice_items"
}

// Coupon represents a discount coupon
type Coupon struct {
	ID              uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Code            string     `json:"code" gorm:"uniqueIndex;not null;size:50"`
	Description     string     `json:"description" gorm:"size:500"`
	DiscountType    string     `json:"discount_type" gorm:"size:20"` // percentage, fixed
	DiscountValue   float64    `json:"discount_value" gorm:"type:decimal(10,2)"`
	MaxUses         int        `json:"max_uses" gorm:"default:0"` // 0 = unlimited
	UsedCount       int        `json:"used_count" gorm:"default:0"`
	MaxUsesPerUser  int        `json:"max_uses_per_user" gorm:"default:1"`
	MinOrderAmount  float64    `json:"min_order_amount" gorm:"type:decimal(10,2);default:0"`
	ApplicablePackages []uuid.UUID `json:"applicable_packages" gorm:"type:jsonb"`
	StartsAt        *time.Time `json:"starts_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	IsActive        bool       `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName returns the table name for Coupon
func (Coupon) TableName() string {
	return "coupons"
}

// IsValid checks if the coupon is currently valid
func (c *Coupon) IsValid() bool {
	now := time.Now()
	if !c.IsActive {
		return false
	}
	if c.MaxUses > 0 && c.UsedCount >= c.MaxUses {
		return false
	}
	if c.StartsAt != nil && now.Before(*c.StartsAt) {
		return false
	}
	if c.ExpiresAt != nil && now.After(*c.ExpiresAt) {
		return false
	}
	return true
}
