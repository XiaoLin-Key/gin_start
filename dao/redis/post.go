package redis

import (
	"gin_start/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// CreatePost 在 Redis 中初始化帖子数据
func CreatePost(postID, CommunityID int64) (err error) {
	now := float64(time.Now().Unix())

	// ✅ 直接接收两个返回值：cmds 和 err
	_, err = Rdb.TxPipelined(Ctx, func(pipe redis.Pipeliner) error {
		pipe.ZAdd(Ctx, getRedisKey(KeyPostTimeZSet), redis.Z{
			Score:  now,
			Member: postID,
		})
		pipe.ZAdd(Ctx, getRedisKey(KeyPostScoreZSet), redis.Z{
			Score:  now,
			Member: postID,
		})
		pipe.SAdd(Ctx, getRedisKey(KeyCommunitySetPF+strconv.Itoa(int(CommunityID))), postID)
		return nil
	})

	// err 会被自动返回（因为你在函数签名里定义了 (err error)）
	return
}

func getIDsFormKey(key string, page, size int64) (ids []string, err error) {
	//获取分页参数
	start := (page - 1) * size
	end := start + size - 1
	//执行查询
	ids, err = Rdb.ZRevRange(Ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func GetPostIDsInOrder(p *models.ParamPostList) (postIDs []string, err error) {
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	return getIDsFormKey(key, p.Page, p.Size)
}

func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	keys := make([]string, 0, len(ids))
	cmds := make([]*redis.IntCmd, 0, len(ids))
	for _, id := range ids {
		keys = append(keys, getRedisKey(KeyPostVotedZSetPF+id))
	}
	_, err = Rdb.Pipelined(Ctx, func(pipe redis.Pipeliner) error {
		for _, key := range keys {
			cmds = append(cmds, pipe.ZCount(Ctx, key, "1", "1"))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	for _, cmd := range cmds {
		data = append(data, cmd.Val())
	}
	return
}

func GetPostCommunityIDsInOrder(p *models.ParamPostList) (postIDs []string, err error) {
	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//zinterstore 合并多个有序集合
	//利用缓存，减少执行次数
	cKey := getRedisKey(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if Rdb.Exists(Ctx, key).Val() == 0 {
		//不存在，执行合并操作
		pipe := Rdb.Pipeline()
		_, err = pipe.ZInterStore(Ctx, key, &redis.ZStore{
			Keys:      []string{cKey, key}, // 所有的源 Key 放在这里
			Aggregate: "MAX",
		}).Result()
		if err != nil {
			return nil, err
		}
		//设置过期时间
		pipe.Expire(Ctx, key, 60*time.Second)
		_, err = pipe.Exec(Ctx)
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
}
