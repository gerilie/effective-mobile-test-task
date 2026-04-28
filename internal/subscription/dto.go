package subscription

const dtoDateLayout = "01-2006"

type (
	SubReq struct {
		ServiceName string  `json:"service_name"       validate:"required"`
		Price       int     `json:"price"              validate:"required,gt=0"`
		UserID      string  `json:"user_id"            validate:"required,uuid4"`
		StartDate   string  `json:"start_date"         validate:"required,datetime=01-2006"`
		EndDate     *string `json:"end_date,omitempty" validate:"omitempty,datetime=01-2006"`
	}
	SubResp struct {
		ID          int     `json:"id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"`
		UserID      string  `json:"user_id"`
		StartDate   string  `json:"start_date"`
		EndDate     *string `json:"end_date,omitempty"`
	}
)

type (
	SubSumReq struct {
		ServiceName string
		UserID      string `validate:"omitempty,uuid4"`
		StartDate   string `validate:"required,datetime=01-2006"`
		EndDate     string `validate:"required,datetime=01-2006"`
	}
	SubSumResp struct {
		TotalPrice int `json:"total_price"`
	}
)

type (
	SubListReq struct {
		ServiceName string
		UserID      string `validate:"omitempty,uuid4"`
		Page        int    `validate:"required,gte=1"`
		Limit       int    `validate:"required,min=1,max=100"`
	}
	SubListResp []SubResp
)

type (
	UpdateSubReq struct {
		ID          int
		ServiceName *string `json:"service_name,omitempty"`
		Price       *int    `json:"price,omitempty"        validate:"omitempty,gt=0"`
		UserID      *string `json:"user_id,omitempty"      validate:"omitempty,uuid4"`
		StartDate   *string `json:"start_date,omitempty"   validate:"omitempty,datetime=01-2006"`
		EndDate     *string `json:"end_date,omitempty"     validate:"required_with=StartDate,omitempty,datetime=01-2006"`
	}
)
