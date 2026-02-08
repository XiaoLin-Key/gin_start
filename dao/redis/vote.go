package redis

import (
	"errors"
	"gin_start/logger"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	oneWeekInSeconds = 7 * 24 * 60 * 60
	scorePreVote     = 432
)

var (
	ErrPostExpired = errors.New("帖子已过期")
)

func VoteForPost(userID, postID string, v float64) (err error) {
	// 1. 取发帖时间（注意是 TimeZSet）
	postTime, err := Rdb.ZScore(Ctx, getRedisKey(KeyPostTimeZSet), postID).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Lg.Error("帖子不存在或已失效", zap.Error(err))
			return errors.New("该帖子不存在或已失效") // 更加明确的错误
		}
		return err
	}
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		return ErrPostExpired
	}

	// 2. 取用户之前的投票纪录
	oldVote := Rdb.ZScore(Ctx, getRedisKey(KeyPostVotedZSetPF+postID), userID).Val()

	// 3. 计算差值逻辑
	var diff float64
	if v == oldVote {
		diff = 0 - oldVote
		v = 0
	} else {
		diff = v - oldVote
	}

	scoreDiff := diff * scorePreVote

	// 4. 关键：使用事务包起来
	_, err = Rdb.TxPipelined(Ctx, func(pipe redis.Pipeliner) error {
		// 更新分数
		if scoreDiff != 0 {
			pipe.ZIncrBy(Ctx, getRedisKey(KeyPostScoreZSet), scoreDiff, postID)
		}
		// 更新用户投票记录
		pipe.ZAdd(Ctx, getRedisKey(KeyPostVotedZSetPF+postID), redis.Z{
			Score:  v,
			Member: userID,
		})
		return nil
	})
	return err
}
