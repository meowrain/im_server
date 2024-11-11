// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package types

type UserInfoRequest struct {
	UserName string `json:"username"`
}

type UserInfoResponse struct {
	UserID         uint   `json:"userID"`
	Nickname       string `json:"nickname"`       // 用户名
	Role           int8   `json:"role"`           // 角色 1 管理员  2 普通用户
	RegisterSource string `json:"registerSource"` // 注册来源
	Abstract       string `json:"abstract"`       // 简介
	Avatar         string `json:"avatar"`         // 头像
}