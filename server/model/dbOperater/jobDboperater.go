package dbOperater

import (
	"demoApp/server/model/dbModel"
	"demoApp/server/model/httpModel"
	"github.com/jinzhu/gorm"
	"goframework/orm"
)

type JobDbOperator struct {
	orm *gorm.DB
}

func (j *JobDbOperator) JobsKind(kind string, offset, limit int) (res []httpModel.HttpJobListModel) {

	res = []httpModel.HttpJobListModel{}
	// TODO 职位分类表
	// 获取职位id
	// 查询具体职位 和 公司
	_ = orm.DB.Model(&dbModel.CompuseJobs{}).Joins("left join  company on company.id = compuse_jobs.company_id").
		Select("compuse_jobs.id as job_id, compuse_jobs.icon_url, compuse_jobs.type, compuse_jobs.name as job_name," +
			"compuse_jobs.location_city as address, compuse_jobs.education as degree, " +
			"compuse_jobs.review_counts as review_count, compuse_jobs.created_time, " +
			"company.name as company_name ").
		Offset(offset).Limit(limit).Scan(&res)

	return res
}

func NewJobDbOperator() *JobDbOperator {

	return &JobDbOperator{
		orm: orm.DB,
	}
}
