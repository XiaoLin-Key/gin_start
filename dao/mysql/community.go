package mysql

import (
	"database/sql"
	"gin_start/logger"
	"gin_start/models"
)

func GetCommunityList() (communityList []*models.Community, err error) {
	//查找所以社区并返回
	sqlStr := "select community_id,community_name from community"
	if err = db.Raw(sqlStr).Scan(&communityList).Error; err != nil {
		if err == sql.ErrNoRows {
			logger.Lg.Warn("there is no community")
			err = nil
		}
	}
	return
}

func GetCommunityDetailByID(id int64) (community *models.CommunityDetail, err error) {
	//查找某个社区并返回
	community = new(models.CommunityDetail)
	sqlStr := "select community_id,community_name,introduction,create_time from community where community_id = ?"
	res := db.Raw(sqlStr, id).Scan(community)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrInvalidID
	}
	return
}
