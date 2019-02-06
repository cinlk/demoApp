package dbOperater

import (
	"demoApp/server/model/dbModel"
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

func NewAppDBoperator() *AppDBoperator {

	return &AppDBoperator{
		orm: orm.DB,
	}
}
