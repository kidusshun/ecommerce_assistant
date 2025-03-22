package coupons

type service struct {
	store CouponStore
}

func NewService(store CouponStore) *service {
	return &service{
		store: store,
	}
}


func (s *service) CreateCoupon(request CreateCouponRequest) (*Coupon, error) {

	createdCoupon, err := s.store.AddCoupon(request.CouponCode, request.DiscountPercentage, request.ExpiryDate, request.Description)
	if err != nil {
		return nil, err
	}

	return createdCoupon, nil
}