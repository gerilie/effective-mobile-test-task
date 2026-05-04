package subscription

import "time"

// sub represents subscription domain model.
type sub struct {
	id          int
	serviceName string
	price       int
	userID      string
	startDate   time.Time
	endDate     *time.Time
}

// subSum represents subscription summa domain model.
type subSum struct {
	serviceName string
	userID      string
	startDate   time.Time
	endDate     time.Time
	totalPrice  int
}

// subList represents subscription list domain model.
type subList struct {
	serviceName string
	userID      string
	page        int
	limit       int
	offset      int
}

// updateSub represents subscription update domain model.
type updateSub struct {
	id          int
	serviceName *string
	price       *int
	userID      *string
	startDate   *time.Time
	endDate     *time.Time
}
