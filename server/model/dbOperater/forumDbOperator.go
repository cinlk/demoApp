package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"goframework/orm"
	"goframework/utils"
)

const (
	hotest    = "hottest"
	interview = "interview"
	recommand = "recommand"
	offer     = "offer"
	help      = "help"
)

type ForumDboperator struct {
	orm *gorm.DB
}

func (f *ForumDboperator) Articles(t string, offset, limit int, userId string) ([]httpModel.HttpForumHttpModel, error) {
	// 如果是热门数据 单独处理
	var res []httpModel.HttpForumHttpModel
	var err error
	if t == hotest {

		err = f.orm.Model(&dbModel.ForumArticle{}).
			Joins("inner join forum_hotest_article on forum_hotest_article.uuid =  forum_article.uuid").
			Joins("inner join \"user\" on \"user\".uuid =  forum_article.user_id").
			Select("forum_article.uuid, forum_article.title, forum_article.content, forum_article.user_id,"+
				"forum_article.read_count as read, forum_article.type as kind, forum_article.created_at as created_time,  "+
				"\"user\".name as user_name, \"user\".user_icon as user_icon").
			Where("forum_article.validate = ?", true).
			Order("forum_hotest_article.created_at desc").Offset(offset).Limit(limit).
			Scan(&res).Error
		if err != nil {
			return nil, err
		}

	} else {
		err = f.orm.Model(&dbModel.ForumArticle{}).
			Joins("inner join \"user\" on \"user\".uuid =  forum_article.user_id").
			Select("forum_article.uuid, forum_article.title, forum_article.content, forum_article.user_id,"+
				"forum_article.read_count as read, forum_article.type as kind, forum_article.created_at as created_time,  "+
				"\"user\".name as user_name, \"user\".user_icon as user_icon").
			Where("forum_article.validate = ? and forum_article.type = ?", true, t).
			Order("forum_article.created_at desc").Offset(offset).Limit(limit).
			Scan(&res).Error
		if err != nil {
			return nil, err
		}
	}

	// 统计点赞和回复数据， 根据userid 检查该用户是否点赞和收藏该帖子
	for i := 0; i < len(res); i++ {

		_ = f.orm.Model(&dbModel.UserLikePost{}).Where("post_uuid = ?", res[i].Uuid).Count(&res[i].ThumbUp)
		_ = f.orm.Model(&dbModel.ReplyForumPost{}).Where("post_uuid = ?", res[i].Uuid).Count(&res[i].Reply)
		var like int
		_ = f.orm.Model(&dbModel.UserLikePost{}).Where("post_uuid = ? and user_id = ?", res[i].Uuid, userId).Count(&like)
		if like != 0 {
			res[i].IsLike = true
		}
		var collected int
		_ = f.orm.Model(&dbModel.UserCollectedPost{}).Where("post_uuid = ? and user_id = ?", res[i].Uuid, userId).Count(&collected)
		if collected != 0 {
			res[i].IsCollected = true
		}

	}

	return res, nil

}

func (f *ForumDboperator) NewArticle(title, content, t, userId string) (string, error) {

	var uid = utils.GetUUID()

	err := f.orm.Create(&dbModel.ForumArticle{
		Uuid:    uid,
		Title:   title,
		UserId:  userId,
		Content: content,
		//ThumbUpCount: 0,
		ReadCount: 1,
		//ReplayCount:  0,
		Type: t,
	}).Error
	return uid, err
}

func (f *ForumDboperator) DeletePostBy(postId, userId string) error {
	session := f.orm.Begin()
	err := session.Where("uuid = ? and user_id = ?", postId, userId).Delete(&dbModel.ForumArticle{}).Error
	if err != nil {
		session.Rollback()
		return err
	}
	err = session.Where("uuid = ?", postId).Delete(&dbModel.ForumHotestArticle{}).Error
	if err != nil {
		session.Rollback()
		return err
	}
	// 需要删除其他关联数据?
	//err = session.Where().Delete(&dbModel.ReplyForumPost{})

	session.Commit()
	return nil
}

func (f *ForumDboperator) PostContentInfo(postId string, offset, limit int) ([]httpModel.HttpSubReplyInfo, error) {
	// 查询帖子 的子回复信息
	var res []httpModel.HttpSubReplyInfo

	err := f.orm.Model(&dbModel.ReplyForumPost{}).Where("post_uuid = ?", postId).
		Joins("inner join \"user\" on \"user\".uuid =  reply_forum_post.user_id").
		Select("\"user\".user_icon, \"user\".name as user_name, reply_forum_post.created_at as created_time, " +
			"reply_forum_post.content, reply_forum_post.reply_id").
		Offset(offset).Limit(limit).Order("reply_forum_post.created_at").
		Scan(&res).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(res); i++ {
		_ = f.orm.Model(&dbModel.UserLikeReply{}).Where("reply_id = ?", res[i].ReplyId).Count(&res[i].LikeCount)
		_ = f.orm.Model(&dbModel.SecondReplyPost{}).Where("reply_id = ?", res[i].ReplyId).Count(&res[i].ReplyCount)
	}

	return res, nil

}

func (f *ForumDboperator) PostReadCount(postId string) error {

	return f.orm.Model(&dbModel.ForumArticle{}).Where("uuid = ?", postId).Update("read_count",
		gorm.Expr("read_count + ?", 1)).Error
}

func (f *ForumDboperator) UserLikePost(userId, postId string, b bool) error {
	session := f.orm.Begin()
	// 记录用户点赞的帖子
	if b {
		// 发送通知给帖子主人
		//leancloud.LeanCloudSendUserNotify()
		err := session.Where("post_uuid = ? and user_id = ?", postId, userId).FirstOrCreate(&dbModel.UserLikePost{
			UserId:   userId,
			PostUuid: postId,
		}).Error
		if err != nil {
			session.Rollback()
			return err
		}

	} else {
		err := session.Unscoped().Delete(dbModel.UserLikePost{}, "user_id = ? and post_uuid =?", userId, postId).Error
		if err != nil {
			session.Rollback()
			return err
		}
	}
	session.Commit()
	return nil
}

func (f *ForumDboperator) UserCollectedPost(userId, postId string, b bool) error {
	// 记录用户收藏的帖子
	session := f.orm.Begin()
	if b {

		err := session.FirstOrCreate(&dbModel.UserCollectedPost{
			UserId:   userId,
			PostUuid: postId,
		}).Error
		if err != nil {
			session.Rollback()
			return err
		}

	} else {
		err := session.Unscoped().Delete(dbModel.UserCollectedPost{}, "user_id = ? and post_uuid =?", userId, postId).Error
		if err != nil {
			session.Rollback()
			return err
		}
	}
	session.Commit()
	return nil

}

// 记录回复帖子的内容, 并通知用户
func (f *ForumDboperator) RecordUserReplyPost(userId, postId, content string) (string, error) {

	var rid = utils.GetUUID()
	session := f.orm.Begin()
	err := session.Create(&dbModel.ReplyForumPost{
		PostUuid: postId,
		UserId:   userId,
		ReplyId:  rid,
		Content:  content,
	}).Error
	if err != nil {
		session.Rollback()
		return "", err
	}
	// 发送通知 TODO
	// _ = leancloud.LeanCloudSendUserNotify()

	session.Commit()

	return rid, nil
}

func NewForumDboperator() *ForumDboperator {
	return &ForumDboperator{
		orm: orm.DB,
	}
}
