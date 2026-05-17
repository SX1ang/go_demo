package dto

import "time"

type CreatePostReq struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CommunityID int64  `json:"community_id" binding:"required"`
}

type GetPostDetailReq struct {
	PostID int64 `uri:"id,string" binding:"required"`
}

type PostDetail struct {
	ID         int64     `json:"id,string"` // 帖子id
	Title      string    `json:"title"`     // 标题
	Content    string    `json:"content"`   // 内容
	Status     int32     `json:"status"`    // 帖子状态
	StatusText string    `json:"status_text"`
	CreateTime time.Time `json:"create_time"` // 创建时间
	UpdateTime time.Time `json:"update_time"` // 更新时间
	VoteCount  int64     `json:"vote_count"`  // 帖子点赞数

	AuthorName string          `json:"author_name"`
	Community  CommunityDetail `json:"community"`
}

type GetPostDetailRes struct {
	Post *PostDetail `json:"post"`
}

type GetPostListReq struct {
	Page  int    `form:"page" binding:"required,min=1"`
	Size  int    `form:"size" binding:"required,min=1,max=100"`
	Order string `form:"order" binding:"omitempty,oneof=new hot"`
}

type GetPostListRes struct {
	Posts []*PostDetail `json:"posts"`
}
