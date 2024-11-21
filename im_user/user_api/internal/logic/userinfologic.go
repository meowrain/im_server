package logic

import (
	"context"
	"encoding/json"
	"errors"

	"im_server/im_user/user_models"
	"im_server/im_user/user_rpc/types/user_rpc"

	"gorm.io/gorm"

	"im_server/im_user/user_api/internal/svc"
	"im_server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户信息获取接口
func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	var user user_models.UserModel
	err = l.svcCtx.DB.Select("id").Where("nickname = ?", req.UserName).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Infof("用户不存在: %s", req.UserName)
			return nil, errors.New("用户不存在")
		} else {
			logx.Errorf("查询失败: %v", err)
			return nil, err
		}
	}

	// 提取用户信息
	res, err := l.svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
		UserId: uint32(user.ID),
	})
	if err != nil {
		logx.Errorf("UserRpc 调用失败: %v", err)
		return nil, err
	}

	// logx.Infof("UserRpc 返回数据: %+v", res)

	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		logx.Errorf("JSON 解析失败: %v", err)
		return nil, errors.New("数据错误")
	}

	// logx.Infof("解析后的用户数据: %+v", user)
	if user.UserConfModel == nil {
		logx.Error("UserConfModel 为空")
	}
	resp = &types.UserInfoResponse{
		UserID:   user.ID,
		Nickname: user.Nickname,
		Abstract: user.Abstract,
		Avatar:   user.Avatar,
	}
	if user.UserConfModel != nil {
		resp.RecallMessage = user.UserConfModel.RecallMessage
		resp.FriendOnline = user.UserConfModel.FriendOnline
		resp.Sound = user.UserConfModel.EnableSound
		resp.SecureLink = user.UserConfModel.SecureLink
		resp.SavePwd = user.UserConfModel.SavePwd
		resp.SearchUser = user.UserConfModel.SearchUser
		resp.Verification = user.UserConfModel.Verification
		if user.UserConfModel.VerificationQuestion != nil {
			resp.VerificationQuestion = &types.VerificationQuestion{
				Problem1: user.UserConfModel.VerificationQuestion.Problem1,
				Problem2: user.UserConfModel.VerificationQuestion.Problem2,
				Problem3: user.UserConfModel.VerificationQuestion.Problem3,
				Answer1:  user.UserConfModel.VerificationQuestion.Answer1,
				Answer2:  user.UserConfModel.VerificationQuestion.Answer2,
				Answer3:  user.UserConfModel.VerificationQuestion.Answer3,
			}
		}
	}
	logx.Infof("最终返回数据: %+v", resp)
	return resp, nil
}