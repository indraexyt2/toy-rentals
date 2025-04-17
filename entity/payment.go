package entity

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

const (
	PaymentTypeRental    = "rental"
	PaymentTypeLateFee   = "late_fee"
	PaymentTypeDamageFee = "damage_fee"
	PaymentTypeCombined  = "combined"
)

const (
	TransactionStatusPending       = "pending"
	TransactionStatusCapture       = "capture"
	TransactionStatusSettlement    = "settlement"
	TransactionStatusDeny          = "deny"
	TransactionStatusCancel        = "cancel"
	TransactionStatusExpire        = "expire"
	TransactionStatusFailure       = "failure"
	TransactionStatusRefund        = "refund"
	TransactionStatusPartialRefund = "partial_refund"
)

type Payment struct {
	BaseEntity
	RentalID          uuid.UUID  `gorm:"type:uuid;not null" json:"rental_id"`
	TransactionID     string     `gorm:"size:100;uniqueIndex" json:"transaction_id"`
	PaymentType       string     `gorm:"size:50;not null;check:payment_type IN ('rental', 'late_fee', 'damage_fee', 'combined')" json:"payment_type"`
	GrossAmount       float64    `gorm:"type:decimal(10,2);not null" json:"gross_amount"`
	SnapToken         string     `gorm:"type:text" json:"snap_token"`
	SnapURL           string     `gorm:"type:text" json:"snap_url"`
	ExpiryTime        *time.Time `json:"expiry_time"`
	TransactionTime   *time.Time `json:"transaction_time"`
	TransactionStatus string     `gorm:"size:50" json:"transaction_status"`
	PaymentMethod     string     `gorm:"size:50" json:"payment_method"`
	VANumber          string     `gorm:"size:100" json:"va_number"`
	FraudStatus       string     `gorm:"size:50" json:"fraud_status"`

	Rental Rental `gorm:"foreignKey:RentalID" json:"-"`
}

func (*Payment) TableName() string {
	return "payments"
}
