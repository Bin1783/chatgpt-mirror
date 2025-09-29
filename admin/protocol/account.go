package protocol

type AddAccountReq struct {
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
	LimitRule    []RateLimitRuleInput `json:"limit_rule" validate:"required"`
	Username     string               `json:"username" validate:"required"`
	Password     string               `json:"password" validate:"required"`
}

type RateLimitRuleInput struct {
	PeriodValue int64  `json:"period_value" validate:"required,min=1"`
	PeriodUnit  string `json:"period_unit" validate:"oneof=second minute hour day"`
	Limit       int64  `json:"limit"     validate:"required,min=1" `
}

type ModelLimitConfig struct {
	EveryMinute int `json:"every_minute"`
	LimitCount  int `json:"limit_count"`
}

type UserRateLimitRule struct {
	UserID int                         `json:"user_id"`
	Rule   map[string]ModelLimitConfig `json:"limit_rule"`
}

type AddUserRateLimitReq struct {
	UserRules []UserRateLimitRule `json:"user_rules"`
}

type GetTokenMessageReq struct {
	Token []string `json:"token"`
}

type GetTokenMessageResp struct {
	AccessToken  string `json:"access_token"`
	SessionToken string `json:"session_token"`
	Email        string `json:"email"`
	PlanType     string `json:"plan_type"`
	ExpireTime   int64  `json:"expire_time"`
	RefreshToken string `json:"refresh_token"`
}
