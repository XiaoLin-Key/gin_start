package logic

import (
	"gin_start/dao/redis"
	"gin_start/models"
	"strconv"
)

//投一票加432分 86400/200

/*投票的几种情况
direction =1,两种情况：
	1.用户之前没有投票，直接投赞成票
	2.用户之前投票了，但是投票方向不同，需要修改投票方向
direction =-1,两种情况：
	1.用户之前没有投票，直接投反对票
	2.用户之前投票了，但是投票方向不同，需要修改投票方向
direction =0,两种情况：
	1.用户之前投赞成票，取消投票
	2.用户之前投反对票，取消投票

投票限制：
每个帖子自发表之日起一个星期内可以投票
	1.到期之后将redis中的投票数据存入mysql
	2.到期之后删除KeyPostVotedZSetPF
*/

func VoteForPost(userID int64, p *models.ParamVoteData) (err error) {
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
}
