package models

// ParamLoginResponse 登录接口返回的数据结构
type ParamLoginResponse struct {
	UserID   int64  `json:"user_id,string"` // string标签防止前端JS丢失大数精度
	Username string `json:"username"`
	Token    string `json:"token"`
}
