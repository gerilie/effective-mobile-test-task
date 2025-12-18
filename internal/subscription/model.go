package subscription

import (
	"errors"
	"time"
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

type SubSumReq struct {
	ServiceName *string `json:"service_name,omitempty"`
	UserID      *string `json:"user_id,omitempty"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     string  `json:"end_date,omitempty"`
}

type SubSumResp struct {
	TotalPrice int `json:"total_price"`
}

var (
	errInvalidUserID = errors.New("invalid user id")

	errEmptyServiceName = errors.New("empty service name")
	errEmptyUserID      = errors.New("empty user id")
	errEmptyServiceUser = errors.New("empty service name and user id. At least one must be used")
	errEmptyPrice       = errors.New("empty price")
	errEmptyStartDate   = errors.New("empty start date")
	errEmptyEndDate     = errors.New("empty end date")
)
