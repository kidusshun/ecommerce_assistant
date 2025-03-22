package coupons

import (
	"database/sql"
	"time"
)

type store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *store {
	return &store{
		db: db,
	}
}


func (s *store) AddCoupon(couponCode string, DiscountPercentage float64, ExpiryDate time.Time, description string) (*Coupon, error) {
	// Add the coupon to the database
	row := s.db.QueryRow("INSERT INTO coupons (coupon_code, discount_percentage, expiration_date, description) VALUES ($1, $2, $3, $4) RETURNING *", couponCode, DiscountPercentage, ExpiryDate, description)

	coupon, err := ScanRowToCoupon(row)

	if err != nil {
		return nil, err
	}

	return coupon, nil
}

func ScanRowToCoupon(row *sql.Row) (*Coupon, error) {
	coupon := new(Coupon)

	err := row.Scan(
		&coupon.ID,
		&coupon.CouponCode,
		&coupon.DiscountPercentage,
		&coupon.ExpiryDate,
		&coupon.IsActive,
		&coupon.CreatedAt,
		&coupon.UpdatedAt,
		&coupon.Description,
	)

	if err != nil {
		return nil, err
	}
	return coupon, nil
}