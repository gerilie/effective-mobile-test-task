package subscription

import (
	"errors"
	"time"
)

var (
	errInvalidUserID = errors.New("invalid user id")

	errEmptyServiceName = errors.New("empty service name")
	errEmptyUserID      = errors.New("empty user id")
	errEmptyPrice       = errors.New("empty price")
	errEmptyStartDate   = errors.New("empty start date")
	errEmptyEndDate     = errors.New("empty end date")
)

type Config struct {
	Host     string
	Port     string
	ReadHTO  time.Duration `mapstructure:"read_header_timeout"`
	ReadTO   time.Duration `mapstructure:"read_timeout"`
	WriteTO  time.Duration `mapstructure:"write_timeout"`
	IdleTO   time.Duration `mapstructure:"idle_timeout"`
	ClientTO time.Duration `mapstructure:"client_timeout"`
}

type SubReq struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type SubResp struct {
	ID          int     `json:"id"`
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
}

type ListSubReq struct {
	page   int
	limit  int
	offset int
}

type ListSubResp []SubResp

type SubSumReq struct {
	serviceName string
	userID      string
	startDate   string
	endDate     string
}

type SubSumResp struct {
	TotalPrice int `json:"total_price"`
}
