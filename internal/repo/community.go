package repo

import (
	"context"
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/infra"
	"gorm.io/gorm"
)

type MysqlCommunityRespository struct {
	db *gorm.DB
}

func NewMysqlCommunityRespository(res *infra.Resources) ICommunityRepo {
	return &MysqlCommunityRespository{db: res.DB}
}

func (m MysqlCommunityRespository) GetCommunityList(ctx context.Context) ([]*model.Community, error) {
	q := query.Use(m.db).WithContext(ctx).Community

	communities, err := q.Find()
	if err != nil {
		return nil, err
	}

	return communities, nil
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
