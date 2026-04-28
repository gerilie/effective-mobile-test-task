package subscription

import "time"

type sub struct {
	id          int
	serviceName string
	price       int
	userID      string
	startDate   time.Time
	endDate     *time.Time
}

type subSum struct {
	serviceName string
	userID      string
	startDate   time.Time
	endDate     time.Time
	totalPrice  int
}

type subList struct {
	serviceName string
	userID      string
	page        int
	limit       int
	offset      int
}

type updateSub struct {
	id          int
	serviceName *string
	price       *int
	userID      *string
	startDate   *time.Time
	endDate     *time.Time
}
