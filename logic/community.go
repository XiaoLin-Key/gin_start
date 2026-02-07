package logic

import (
	"gin_start/dao/mysql"
	"gin_start/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	//查找所以社区并返回
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (data *models.CommunityDetail, err error) {
	//查找某个社区详情并返回
	return mysql.GetCommunityDetailByID(id)
}
