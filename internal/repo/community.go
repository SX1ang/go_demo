package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/infra"
	"gorm.io/gorm"
)

type CommunityBasic struct {
	CommunityID   uint   `gorm:"column:community_id"`
	CommunityName string `gorm:"column:community_name"`
}

type MysqlCommunityRespository struct {
	db *gorm.DB
}

func NewMysqlCommunityRespository(res *infra.Resources) ICommunityRepo {
	return &MysqlCommunityRespository{db: res.DB}
}

func (m MysqlCommunityRespository) GetCommunityList(ctx context.Context) ([]*CommunityBasic, error) {
	q := query.Use(m.db)

	var result []*CommunityBasic
	err := q.WithContext(ctx).
		Community.
		Select(q.Community.CommunityID, q.Community.CommunityName).
		Scan(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m MysqlCommunityRespository) GetCommunityDetail(ctx context.Context, communityId int64) (*model.Community, error) {
	q := query.Use(m.db)

	community, err := q.Community.
		WithContext(ctx).
		Where(q.Community.CommunityID.Eq(int32(communityId))).
		First()

	if err != nil {
		return nil, err
	}

	return community, nil
}

func (m MysqlCommunityRespository) AddCommunity(ctx context.Context, community *model.Community) error {
	err := query.Use(m.db).WithContext(ctx).Community.Create(community)
	if err != nil {
		return err
	}
	return nil
}
