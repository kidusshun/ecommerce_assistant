package coupons

import (
	"time"

	"github.com/google/uuid"
)

type CouponService interface {
	CreateCoupon(request CreateCouponRequest) (*Coupon, error)
}

type CouponStore interface {
	AddCoupon(couponCode string, DiscountPercentage float64, ExpiryDate time.Time, description string) (*Coupon, error)
}


type Coupon struct {
	ID uuid.UUID `json:"id"`
	CouponCode string `json:"coupon_code"`
	DiscountPercentage float64 `json:"discount_percentage"`
	Description string `json:"description"`
	ExpiryDate time.Time `json:"expiration_date"`
	IsActive bool `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateCouponRequest struct {
	CouponCode string `json:"coupon_code"`
	DiscountPercentage float64 `json:"discount_percentage"`
	ExpiryDate time.Time `json:"expiration_date"`
	Description string `json:"description"`
}