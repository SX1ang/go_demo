package service

import (
	"demoProject/internal/api/dto"
	"demoProject/internal/repo"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommunityService struct {
	communityRepo repo.ICommunityRepo
}

func NewCommunityService(communityRepo repo.ICommunityRepo) ICommunityService {
	return &CommunityService{
		communityRepo: communityRepo,
	}
}

func (comm *CommunityService) GetCommunityList(ctx *gin.Context) (*dto.GetCommunityListRes, error) {
	commList, err := comm.communityRepo.GetCommunityList(ctx)
	if err != nil {
		zap.L().Error("GetCommunityList failed", zap.Error(err))
		return nil, err
	}

	communities := make([]*dto.CommunityItem, 0, len(commList))
	for _, comm := range commList {
		communities = append(communities, &dto.CommunityItem{
			ID:   int64(comm.CommunityID),
			Name: comm.CommunityName,
		})
	}

	return &dto.GetCommunityListRes{
		Communities: communities,
	}, nil
}

func (comm *CommunityService) GetCommunityDetail(ctx *gin.Context, req *dto.GetCommunityDetailReq) (*dto.GetCommunityDetailRes, error) {
	commDetail, err := comm.communityRepo.GetCommunityDetail(ctx, req.CommunityId)
	if err != nil {
		zap.L().Error("GetCommunityDetail failed", zap.Error(err))
		return nil, err
	}

	return &dto.GetCommunityDetailRes{
		Community: &dto.CommunityDetail{
			ID:           int64(commDetail.CommunityID),
			Name:         commDetail.CommunityName,
			Introduction: commDetail.Introduction,
			CreateTime:   commDetail.CreateTime,
			UpdateTime:   commDetail.UpdateTime,
		},
	}, nil
}
