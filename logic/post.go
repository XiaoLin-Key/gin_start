package logic

import (
	"gin_start/dao/mysql"
	"gin_start/dao/redis"
	"gin_start/logger"
	"gin_start/models"
	"gin_start/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//生成post id
	p.ID = int64(snowflake.GenID())
	//保存数据库
	if err := mysql.CreatePost(p); err != nil {
		return err
	}
	if err = redis.CreatePost(p.ID, p.CommunityID); err != nil {
		logger.Lg.Error("redis.CreatePost failed", zap.Int64("post_id", p.ID), zap.Error(err))
	}
	return nil
}

func GetPostById(id int64) (data *models.PostDetailVO, err error) {
	//获取帖子详情
	post, err := mysql.GetPostById(id)
	if err != nil {
		logger.Lg.Error("mysql.GetPostById() failed, err: %v", zap.Int64("post_id", id), zap.Error(err))
		return nil, err
	}
	//获取作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		logger.Lg.Error("mysql.GetUserByID() failed, err: %v", zap.Int64("user_id", post.AuthorID), zap.Error(err))
		return nil, err
	}
	//获取社区信息
	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		logger.Lg.Error("mysql.GetCommunityDetailByID() failed, err: %v", zap.Int64("community_id", post.CommunityID), zap.Error(err))
		return nil, err
	}

	//组织返回数据
	data = &models.PostDetailVO{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return
}

func GetPostList(page, size int64) (data []*models.PostDetailVO, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.PostDetailVO, 0, len(posts))
	for _, post := range posts {
		//获取作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			logger.Lg.Error("mysql.GetUserByID() failed, err: %v", zap.Int64("user_id", post.AuthorID), zap.Error(err))
			continue
		}
		//获取社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			logger.Lg.Error("mysql.GetCommunityDetailByID() failed, err: %v", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.PostDetailVO{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.PostDetailVO, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	data = make([]*models.PostDetailVO, 0, len(posts))
	for idx, post := range posts {
		//获取作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			logger.Lg.Error("mysql.GetUserByID() failed, err: %v", zap.Int64("user_id", post.AuthorID), zap.Error(err))
			continue
		}
		//获取社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			logger.Lg.Error("mysql.GetCommunityDetailByID() failed, err: %v", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.PostDetailVO{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList2(p *models.ParamPostList) (data []*models.PostDetailVO, err error) {
	ids, err := redis.GetPostCommunityIDsInOrder(p)
	if err != nil {
		return nil, err
	}
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return nil, err
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}
	data = make([]*models.PostDetailVO, 0, len(posts))
	for idx, post := range posts {
		//获取作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			logger.Lg.Error("mysql.GetUserByID() failed, err: %v", zap.Int64("user_id", post.AuthorID), zap.Error(err))
			continue
		}
		//获取社区信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			logger.Lg.Error("mysql.GetCommunityDetailByID() failed, err: %v", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postDetail := &models.PostDetailVO{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (data []*models.PostDetailVO, err error) {
	if p.CommunityID == 0 {
		return GetPostList2(p)
	}
	return GetCommunityPostList2(p)
}
