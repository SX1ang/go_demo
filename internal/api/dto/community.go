package dto

import (
	"time"
)

type CommunityItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type GetCommunityListRes struct {
	Communities []*CommunityItem `json:"communities"`
}

type GetCommunityDetailReq struct {
	CommunityId int64 `uri:"id" binding:"required"`
}

type CommunityDetail struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Introduction string    `json:"introduction"`
	CreateTime   time.Time `json:"create_time"`
	UpdateTime   time.Time `json:"update_time"`
}

type GetCommunityDetailRes struct {
	Community *CommunityDetail `json:"community"`
}
