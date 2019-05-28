package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"goframework/orm"
)

type AppDBoperator struct {
	orm *gorm.DB
}

func (a *AppDBoperator) AppGuidanceItems() (res []dbModel.AppGuidanceItem) {

	_ = a.orm.Model(&dbModel.AppGuidanceItem{}).Find(&res)
	return

}

func (a *AppDBoperator) News(t string, offset int) (res []dbModel.NewsModel) {

	_ = a.orm.Model(&dbModel.NewsModel{}).Offset(offset).Limit(10).Find(&res)

	return
}

func (a *AppDBoperator) UserRelatedPostGroup(userId string) ([]httpModel.UserPostGroups, error) {

	var res []httpModel.UserPostGroups


	err := a.orm.Model(&dbModel.UserCollectedGroup{}).Where("user_id = ?", userId).
		Select("group_name as name, id as group_id ").Order("created_at desc").
		Scan(&res).Error

	//for _, i := range  name{
	//	res.Name = append(res.Name, i.Name)
	//}
	return res, err
}


func (a *AppDBoperator) UserTalkDefaultMessage(userId string) (string, error){

	var target dbModel.DefaultFirstMessage
	err := a.orm.Model(&target).Where("user_id = ?", userId).First(&target).Error
	if err != nil{
		return "", err
	}

	return target.Messages[target.DefaultNum], nil
}

func NewAppDBoperator() *AppDBoperator {

	return &AppDBoperator{
		orm: orm.DB,
	}
}
