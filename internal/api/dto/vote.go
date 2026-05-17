package dto

type VotePostReq struct {
	PostID int64 `json:"post_id,string" binding:"required"`
	Vote   int   `json:"vote" binding:"oneof=-1 0 1"`
	// 1: 赞成 0: 取消 -1: 反对
}
