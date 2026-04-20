package repo

import (
	"demoProject/dal/model"
	"demoProject/dal/query"
	"demoProject/internal/infra"
	"errors"

	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlSessionRespository struct {
	db *gorm.DB
}

func NewMysqlSessionRespository(res *infra.Resources) ISessionRepo {
	return &MysqlSessionRespository{
		db: res.DB,
	}
}

func (m *MysqlSessionRespository) SaveSession(ctx context.Context, session *model.Session) error {
	err := query.Use(m.db).
		WithContext(ctx).
		Session.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"session_id",
				"refresh_token",
				"is_revoked",
				"created_at",
				"expires_at",
			}),
		}).
		Create(session)
	if err != nil {
		return err
	}

	return nil
}

func (m *MysqlSessionRespository) GetSession(ctx context.Context, sessionId string) (*model.Session, error) {
	q := query.Use(m.db)
	session, err := q.WithContext(ctx).Session.Where(q.Session.SessionID.Eq(sessionId)).First()
	if err == nil {
		return session, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err

}

func (m *MysqlSessionRespository) DeleteSession(ctx context.Context, userId int64) error {
	q := query.Use(m.db)
	result, err := q.WithContext(ctx).Session.Where(q.Session.UserID.Eq(userId)).Delete()

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}

func (m *MysqlSessionRespository) RevokeRefreshToken(ctx context.Context, userId int64) error {
	q := query.Use(m.db)
	result, err := q.WithContext(ctx).Session.
		Where(
			q.Session.UserID.Eq(userId),
			q.Session.IsRevoked.Is(false),
		).
		UpdateSimple(
			q.Session.IsRevoked.Value(true),
		)

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return errors.New("session not found")
	}

	return nil
}
