package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"goframework/orm"
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

func (f *ForumDboperator) Articles(t string, offset, limit int) ([]httpModel.HttpForumHttpModel, error) {
	// 如果是热门数据 单独处理
	var res []httpModel.HttpForumHttpModel
	var err error
	if t == hotest {

		err = f.orm.Model(&dbModel.ForumArticle{}).
			Joins("inner join forum_hotest_article on forum_hotest_article.uuid =  forum_article.uuid").
			Joins("inner join \"user\" on \"user\".uuid =  forum_article.user_id").
			Select("forum_article.uuid, forum_article.title, forum_article.user_id,"+
				"forum_article.thumb_up_count as thumb_up, forum_article.replay_count as reply, "+
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
			Select("forum_article.uuid, forum_article.title, forum_article.user_id,"+
				"forum_article.thumb_up_count as thumb_up, forum_article.replay_count as reply, "+
				"forum_article.read_count as read, forum_article.type as kind, forum_article.created_at as created_time,  "+
				"\"user\".name as user_name, \"user\".user_icon as user_icon").
			Where("forum_article.validate = ? and forum_article.type = ?", true, t).
			Order("forum_article.created_at desc").Offset(offset).Limit(limit).
			Scan(&res).Error
		if err != nil {
			return nil, err
		}
	}

	return res, nil

}

func NewForumDboperator() *ForumDboperator {
	return &ForumDboperator{
		orm: orm.DB,
	}
}
