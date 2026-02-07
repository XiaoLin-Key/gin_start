package mysql

import "gin_start/models"

func CreatePost(p *models.Post) (err error) {
	sqlStr := "insert into post (post_id, author_id, community_id, status, title, content) values (?, ?, ?, ?, ?, ?)"
	res := db.Exec(sqlStr, p.ID, p.AuthorID, p.CommunityID, p.Status, p.Title, p.Content)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 1 {
		return ErrorInsertFailed
	}
	return
}

func GetPostById(id int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := "select post_id, author_id, community_id, create_time, title, content from post where post_id = ?"
	res := db.Raw(sqlStr, id).Scan(post)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected != 1 {
		return nil, ErrInvalidID
	}
	return
}

func GetPostList(page, size int64) (postList []*models.Post, err error) {
	postList = make([]*models.Post, 0, size)
	sqlStr := "select post_id, author_id, community_id, create_time, title, content from post limit ?, ?"
	res := db.Raw(sqlStr, (page-1)*size, size).Scan(&postList)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, ErrInvalidID
	}
	return
}
