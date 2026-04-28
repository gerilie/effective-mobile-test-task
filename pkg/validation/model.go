package validation

type Errors map[string]string

type Resp struct {
	Fields Errors `json:"fields"`
}
