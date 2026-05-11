package dto

import "demoProject/dal/model"

type GetCommunityListRes struct {
	Communities []*model.Community `json:"communities"`
}

type GetCommunityDetailReq struct {
	CommunityId int64 `uri:"id" binding:"required"`
}

type GetCommunityDetailRes struct {
	Community *model.Community `json:"community"`
}
