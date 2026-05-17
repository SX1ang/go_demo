package repo

import (
	"context"
	"demoProject/internal/e"
	"demoProject/internal/infra"
	"demoProject/pkg/util"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	oneWeekSeconds = 7 * 24 * 60 * 60
	scorePeerVote  = 432 // 每一票值432分
)

const (
	KeyPostTime  = "demo:post:time"  // 帖子和发布时间
	KeyPostScore = "demo:post:score" // 帖子和帖子分数
)

func KeyPostVoted(postId int64) string {
	return fmt.Sprintf("demo:post:voted:%d", postId) // 每个帖子投票用户和投票类型
}

type VoteRepository struct {
	client *redis.Client
}

func NewVoteRepository(res *infra.Resources) IVoteRepo {
	return &VoteRepository{
		client: res.Redis,
	}
}

func (r *VoteRepository) VotePost(ctx context.Context, userId, postId int64, vote int) error {
	votedKey := KeyPostVoted(postId)
	userIdStr := strconv.FormatInt(userId, 10)
	postIdStr := strconv.FormatInt(postId, 10)

	// 1.判断帖子投票限制，发布超过一周的帖子不支持投票（否则需要为每个帖子记录投票数据）
	postTime := r.client.ZScore(ctx, KeyPostTime, postIdStr).Val()
	if float64(time.Now().Unix())-postTime > oneWeekSeconds {
		return &util.CustomizedErr{
			Code: e.ERROR_VOTE_TIME_EXPIRED,
			Msg:  e.GetMsg(e.ERROR_VOTE_TIME_EXPIRED),
		}
	}

	// 2.更新帖子分数，先检查用户之前的投票记录，计算出要更新的分数
	original_vote := r.client.ZScore(ctx, votedKey, userIdStr).Val()

	// 计算两次投票的差值，决定修改多少分 diff*scorePeerVote
	// 之前投反对票现在改投赞成票 或者 之前投赞成票现在改投反对票 差值都为2，通过vote的值决定是加分还是减分
	// 之前没投票现在改投赞成票 或者 之前没投票现在改投反对票 或者 之前投赞成票/反对票现在取消投票 差值都为1，通过vote的值决定是加分还是减分
	diff := math.Abs(float64(vote) - original_vote)

	// 决定加分还是减分
	var op float64
	if original_vote < float64(vote) {
		op = 1
	} else {
		op = -1
	}

	pipeline := r.client.TxPipeline()

	// 更新分数
	pipeline.ZIncrBy(ctx, KeyPostScore, op*diff*scorePeerVote, postIdStr)

	// 3.记录用户的投票数据
	if vote == 0 {
		// 取消投票，删除记录
		pipeline.ZRem(ctx, votedKey, userIdStr)
	} else {
		// 投票，更新记录（之前投过票，则更新，否则创建）
		pipeline.ZAdd(ctx, votedKey, redis.Z{
			// 键值对
			Score:  float64(vote), // 赞成/反对/取消
			Member: userIdStr,     // 投票的用户ID
		})
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (r *VoteRepository) GetPostsVoteCount(ctx context.Context, ids []int64) ([]int64, error) {
	// 使用pipeline发送多条命令，减少redis连接次数
	pipeline := r.client.TxPipeline()
	for _, id := range ids {
		// 统计赞成票的数量，用于前端展示
		pipeline.ZCount(ctx, KeyPostVoted(id), "1", "1")
	}

	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}

	counts := make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		counts = append(counts, cmder.(*redis.IntCmd).Val())
	}

	return counts, nil
}
