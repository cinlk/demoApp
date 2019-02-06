package dbOperater

import (
	"demoApp/server/model/dbModel"
	"github.com/jinzhu/gorm"
	"goframework/gLog"
	"goframework/orm"
)

type ListDboperater struct {
	orm *gorm.DB
}

func (l *ListDboperater) Banners(limit int) []dbModel.Banners {

	var banner = []dbModel.Banners{}
	_ = orm.DB.Model(&dbModel.Banners{}).Find(&banner).Limit(limit)

	return banner

}

func (l *ListDboperater) LatestNews() []dbModel.LatestNews {

	var news = []dbModel.LatestNews{}
	_ = orm.DB.Model(&dbModel.LatestNews{}).Find(&news)
	return news
}

func (l *ListDboperater) LatestJobCategory(limit int) []dbModel.JobCategory {
	var jobs = []dbModel.JobCategory{}
	_ = orm.DB.Model(&dbModel.JobCategory{}).Find(&jobs).Limit(limit)

	return jobs
}

func (l *ListDboperater) TopJobs() []dbModel.TopJobs {

	var jobs = []dbModel.TopJobs{}
	_ = orm.DB.Model(&dbModel.TopJobs{}).Find(&jobs)
	return jobs
}

func (l *ListDboperater) CarrerTalks(limit int) []dbModel.CarrerTalk {

	var talks = []dbModel.CarrerTalk{}
	_ = orm.DB.Model(&dbModel.CarrerTalk{}).Limit(limit).Find(&talks)
	needDelete := []int{}
	// 查找公司数据
	for i := 0; i < len(talks); i++ {
		err := orm.DB.Model(&dbModel.Company{}).Where("id = ?", talks[i].CompanyID).Find(&talks[i].Company).Error
		if err != nil {
			gLog.LOG_ERROR(err)
			needDelete = append(needDelete, i)
			continue
		}
	}
	// 删除指定数据

	return talks
}

func (l *ListDboperater) OnlineApplyClass() []dbModel.ApplyClassify {

	var c = []dbModel.ApplyClassify{}
	_ = orm.DB.Model(&dbModel.ApplyClassify{}).Find(&c)
	return c
}

func (l *ListDboperater) Jobs(offset, limit int) []dbModel.CompuseJobs {

	var jobs = []dbModel.CompuseJobs{}

	_ = orm.DB.Model(&dbModel.CompuseJobs{}).Offset(offset).Limit(limit).Find(&jobs)

	// 获取关联的数据 公司 和发布职位者
	for i := 0; i < len(jobs); i++ {
		orm.DB.Model(&dbModel.Company{}).Where("id = ?", jobs[i].CompanyID).Find(&jobs[i].Company)
	}
	return jobs
}

func NewListDboperater() *ListDboperater {

	return &ListDboperater{
		orm: orm.DB,
	}

}
